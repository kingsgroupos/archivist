package archivist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/edwingeng/slog"
	"github.com/pkg/errors"
	"gitlab-ee.funplus.io/watcher/watcher/misc"
	"gitlab-ee.funplus.io/watcher/watcher/misc/wtime"
	"go.uber.org/atomic"
)

var (
	ErrUpgradeNeeded = fmt.Errorf("%s", "configuration upgrade needed")
)

type Collection interface {
	CompatibleVersions() []string
	Filename2Conf() map[string]interface{}
	RevRefGraph() map[string][]string
	BindRefs() error
	FixPointers()
}

type Archivist struct {
	slog.Logger

	rootDir        string
	whitelist      []string
	blacklist      []string
	reloadCallback func(newC, oldC Collection) error

	mu      sync.Mutex
	wrapper atomic.Value
}

func NewArchivist(opts ...Option) *Archivist {
	arch := &Archivist{
		Logger: slog.NewConsoleLogger(),
	}
	if workingDir, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		arch.rootDir = filepath.Join(workingDir, "conf")
	}
	for _, opt := range opts {
		opt(arch)
	}
	return arch
}

func (this *Archivist) PatchCollection(coll Collection, overwrites ...Overwrite) (Collection, error) {
	if len(overwrites) == 0 {
		return coll, nil
	}
	if coll.RevRefGraph() == nil {
		return nil, errors.New("RevRefGraph() returns nil")
	}

	this.Info("<archivist> caution: PatchCollection does NOT update collection extension")

	newObj := reflect.New(reflect.TypeOf(coll).Elem()).Interface()
	updateElemViaReflection(newObj, coll)
	newC := newObj.(Collection)
	newC.FixPointers()

	var organized organizedOverwrites
	if err := organized.organize(newC, overwrites); err != nil {
		return nil, err
	}

	for filename := range organized.fileLevel {
		v := newC.Filename2Conf()[filename]
		this.Infof("apply file-level patch to %s", filename)
		vv := reflect.New(reflect.TypeOf(v).Elem()).Interface()
		err := this.loadInmemoryOverwrite(organized, filename, vv)
		if err != nil {
			return nil, err
		}
		updateElemViaReflection(v, vv)
	}

	for filename, patches := range organized.contentLevel {
		v := newC.Filename2Conf()[filename]
		this.Infof("apply content-level patch to %s", filename)
		b, err := json.Marshal(v)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		vv := reflect.New(reflect.TypeOf(v).Elem()).Interface()
		err = json.Unmarshal(b, vv)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		for _, patch := range patches {
			err := applyAdjustmentsImpl(patch, vv)
			if err != nil {
				return nil, err
			}
		}
		updateElemViaReflection(v, vv)
	}

	organizedAffected := organizedOverwrites{
		fileLevel: make(map[string][][]byte),
	}
	affected := organized.findAffected(coll.RevRefGraph())
	for _, filename := range affected {
		v := newC.Filename2Conf()[filename]
		this.Infof("clone %s because it is affected by the patch", filename)
		b, err := json.Marshal(v)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		organizedAffected.fileLevel[filename] = [][]byte{b}
	}
	for filename := range organizedAffected.fileLevel {
		v := newC.Filename2Conf()[filename]
		vv := reflect.New(reflect.TypeOf(v).Elem()).Interface()
		err := this.loadInmemoryOverwrite(organizedAffected, filename, vv)
		if err != nil {
			return nil, err
		}
		updateElemViaReflection(v, vv)
	}

	if err := newC.BindRefs(); err != nil {
		return nil, err
	}
	return newC, nil
}

func updateElemViaReflection(whose1, whose2 interface{}) {
	reflect.ValueOf(whose1).Elem().Set(reflect.ValueOf(whose2).Elem())
}

func (this *Archivist) LoadCollection(newCollection func() interface{}, overwrites ...Overwrite) (Collection, error) {
	g := ConfGroup()
	if g == "" {
		return nil, errors.WithStack(errEnv)
	}
	p := filepath.Join(this.rootDir, g)
	if err := misc.FindDirectory(p); err != nil {
		return nil, errors.WithStack(err)
	}
	subg := ConfSubgroup()
	if subg != "" {
		if err := misc.FindDirectory(filepath.Join(this.rootDir, g, subg)); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	newC, ok := newCollection().(Collection)
	if !ok {
		return nil, errors.New("failed to create a new collection")
	}
	var organized organizedOverwrites
	if err := organized.organize(newC, overwrites); err != nil {
		return nil, err
	}

	shouldIgnore := func(filename string) bool {
		if this.whitelist != nil {
			if misc.IndexStrings(this.whitelist, filename) == -1 {
				return true
			}
		} else if this.blacklist != nil {
			for _, f := range this.blacklist {
				if f == filename {
					return true
				}
			}
		}
		return false
	}

	check1 := func(p1 string, root bool) error {
		err1 := misc.FindFile(p1)
		p2 := strings.TrimSuffix(p1, ".json") + ".js"
		err2 := misc.FindFile(p2)
		if root && err1 != nil && err2 != nil {
			return errors.Errorf("cannot find %s or %s under the root directory",
				filepath.Base(p1), filepath.Base(p2))
		}
		if err1 == nil && err2 == nil {
			return errors.Errorf("%s and %s cannot coexist under the same directory",
				filepath.Base(p1), filepath.Base(p2))
		}
		return nil
	}
	check2 := func(p1, suffix string) error {
		p2 := strings.TrimSuffix(p1, ".json") + suffix
		if err := misc.FindFile(p2); err == nil {
			return errors.Errorf("%s file should not reside under the root directory. file: %s",
				suffix, filepath.Base(p2))
		}
		return nil
	}
	check3 := func(filename, suffix string) error {
		filename = strings.TrimSuffix(filename, ".json") + suffix
		p2 := filepath.Join(this.rootDir, g, subg, filename)
		if err := misc.FindFile(p2); err == nil {
			rel := filepath.Join(g, subg, filename)
			return errors.Errorf("subgroup does not support file-level overwriting. file: %s", rel)
		}
		return nil
	}

	for filename := range newC.Filename2Conf() {
		if shouldIgnore(filename) {
			continue
		}
		p1 := filepath.Join(this.rootDir, filename)
		if err := check1(p1, true); err != nil {
			return nil, err
		}
		if err := check2(p1, ".tweak.json"); err != nil {
			return nil, err
		}
		if err := check2(p1, ".tweak.js"); err != nil {
			return nil, err
		}
		if err := check2(p1, ".local.json"); err != nil {
			return nil, err
		}
		if err := check2(p1, ".local.js"); err != nil {
			return nil, err
		}
		if subg == "" {
			continue
		}
		if err := check3(filename, ".json"); err != nil {
			return nil, err
		}
		if err := check3(filename, ".js"); err != nil {
			return nil, err
		}
	}

	for filename := range newC.Filename2Conf() {
		if shouldIgnore(filename) {
			continue
		}
		p1 := filepath.Join(this.rootDir, g, filename)
		if err := check1(p1, false); err != nil {
			return nil, err
		}
		if subg == "" {
			continue
		}
		p2 := filepath.Join(this.rootDir, g, subg, filename)
		if err := check1(p2, false); err != nil {
			return nil, err
		}
	}

	for filename, v := range newC.Filename2Conf() {
		if shouldIgnore(filename) {
			continue
		}
		var loaded bool
		if organized.fileLevel != nil {
			if _, ok := organized.fileLevel[filename]; ok {
				err := this.loadInmemoryOverwrite(organized, filename, v)
				if err != nil {
					return nil, err
				}
				loaded = true
			}
		}
		if !loaded {
			err := this.loadFile(filename, v)
			if err != nil {
				return nil, err
			}
		}
		if organized.contentLevel != nil {
			if _, ok := organized.contentLevel[filename]; ok {
				this.Infof("start to load %s (in-memory content)...", filename)
				for _, data := range organized.contentLevel[filename] {
					err := applyAdjustmentsImpl(data, v)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}

	if err := newC.BindRefs(); err != nil {
		return nil, err
	}
	if this.reloadCallback != nil {
		var oldC Collection
		if wrapper := this.collectionWrapper(); wrapper != nil {
			oldC = wrapper.Collection
		}
		if err := this.reloadCallback(newC, oldC); err != nil {
			return nil, err
		}
	}

	return newC, nil
}

func (this *Archivist) loadInmemoryOverwrite(organized organizedOverwrites, filename string, v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("<archivist> filename: %s, panic: %+v\n%s", filename, r, debug.Stack())
		}
	}()

	stopwatch := wtime.NewStopwatch()
	err = this.loadInmemoryOverwriteImpl(organized, filename, v)
	if err != nil {
		return err
	}

	dt := int(stopwatch.ElapsedMilliseconds())
	this.Infof("loaded %s (in-memory file). elapsed: %dms", filename, dt)
	return nil
}

func (this *Archivist) loadInmemoryOverwriteImpl(organized organizedOverwrites, filename string, v interface{}) error {
	this.Infof("start to load %s (in-memory file)...", filename)
	n := len(organized.fileLevel[filename])
	if n == 0 {
		panic("impossible")
	}
	err := json.Unmarshal(organized.fileLevel[filename][n-1], v)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (this *Archivist) loadFile(filename string, v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("<archivist> filename: %s, panic: %+v\n%s", filename, r, debug.Stack())
		}
	}()

	stopwatch := wtime.NewStopwatch()
	err = this.loadFileImpl(filename, v)
	if err != nil {
		return err
	}

	dt := int(stopwatch.ElapsedMilliseconds())
	this.Infof("loaded %s. elapsed: %dms", filename, dt)
	return nil
}

func (this *Archivist) loadFileImpl(filename string, v interface{}) error {
	g := ConfGroup()
	rels := []string{
		filepath.Join(".runtime", g, filename),
		filepath.Join(g, filename),
		filepath.Join(".runtime", filename),
		filepath.Join(filename),
	}
	filePaths := []string{
		filepath.Join(this.rootDir, rels[0]),
		filepath.Join(this.rootDir, rels[1]),
		filepath.Join(this.rootDir, rels[2]),
		filepath.Join(this.rootDir, rels[3]),
	}
	var fp, rel string
	for i, filePath := range filePaths {
		if err := misc.FindFile(filePath); err == nil {
			fp, rel = filePath, rels[i]
			break
		}
	}
	if fp == "" {
		return errors.Errorf("cannot find %s anywhere", filename)
	}

	this.Infof("start to load %s...", rel)
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return errors.WithStack(err)
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return errors.WithStack(err)
	}

	err = this.applyAdjustments(g, filename, v, "tweak")
	if err != nil {
		return err
	}
	err = this.applyAdjustments(g, filename, v, "local")
	if err != nil {
		return err
	}

	subg := ConfSubgroup()
	if subg == "" {
		return nil
	}

	subdir := filepath.Join(g, subg)
	err = this.applyAdjustments(subdir, filename, v, "tweak")
	if err != nil {
		return err
	}
	err = this.applyAdjustments(subdir, filename, v, "local")
	if err != nil {
		return err
	}
	return nil
}

func (this *Archivist) applyAdjustments(subdir, filename string, v interface{}, changeType string) error {
	fName := fmt.Sprintf("%s.%s.json", strings.TrimSuffix(filename, ".json"), changeType)
	rels := []string{
		filepath.Join(".runtime", subdir, fName),
		filepath.Join(subdir, fName),
	}
	filePaths := []string{
		filepath.Join(this.rootDir, rels[0]),
		filepath.Join(this.rootDir, rels[1]),
	}
	var fp, rel string
	for i, filePath := range filePaths {
		if err := misc.FindFile(filePath); err == nil {
			fp, rel = filePath, rels[i]
			break
		}
	}
	if fp == "" {
		return nil
	}

	this.Infof("start to load %s...", rel)
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return errors.WithStack(err)
	}

	return applyAdjustmentsImpl(data, v)
}

func applyAdjustmentsImpl(data []byte, v interface{}) error {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Map || t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Map {
		var m2 map[string]interface{}
		if err := json.Unmarshal(data, &m2); err != nil {
			return errors.WithStack(err)
		}
		vv := reflect.ValueOf(v)
		if t.Kind() != reflect.Map {
			vv = vv.Elem()
		}
		var zeroValue reflect.Value
		m3 := make(map[string]struct{}, vv.Len())
		for _, mk := range vv.MapKeys() {
			mkStr := fmt.Sprint(mk.Interface())
			m3[mkStr] = struct{}{}
			if v2, ok := m2[mkStr]; ok {
				j, err := json.Marshal(v2)
				if err != nil {
					return errors.WithStack(err)
				}
				var p interface{}
				if mv := vv.MapIndex(mk); mv != zeroValue {
					if mv.Kind() == reflect.Ptr {
						p = mv.Interface()
					} else if mv.CanAddr() {
						p = mv.Addr().Interface()
					}
				}
				if p != nil {
					if err := json.Unmarshal(j, p); err != nil {
						return errors.WithStack(err)
					}
				} else {
					newElem := reflect.New(vv.Type().Elem())
					if err := json.Unmarshal(j, newElem.Interface()); err != nil {
						return errors.WithStack(err)
					}
					vv.SetMapIndex(mk, newElem.Elem())
				}
			}
		}

		str2KeyValue := str2KeyValueFunc(vv.Type().Key().Kind())
		for mk := range m2 {
			if _, ok := m3[mk]; ok {
				continue
			}
			j, err := json.Marshal(m2[mk])
			if err != nil {
				return errors.WithStack(err)
			}
			newElem := reflect.New(vv.Type().Elem())
			if err := json.Unmarshal(j, newElem.Interface()); err != nil {
				return errors.WithStack(err)
			}
			mkv, err := str2KeyValue(mk)
			if err != nil {
				return errors.WithStack(err)
			}
			vv.SetMapIndex(mkv, newElem.Elem())
		}
	} else {
		if err := json.Unmarshal(data, v); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func str2KeyValueFunc(kind reflect.Kind) func(string) (reflect.Value, error) {
	switch kind {
	case reflect.Bool:
		return func(str string) (reflect.Value, error) {
			v, err := strconv.ParseBool(str)
			return reflect.ValueOf(v), err
		}
	case reflect.Int:
		return func(str string) (reflect.Value, error) {
			v, err := strconv.Atoi(str)
			return reflect.ValueOf(v), err
		}
	case reflect.Int8:
		return func(str string) (reflect.Value, error) {
			v, err := strconv.ParseInt(str, 10, 8)
			return reflect.ValueOf(int8(v)), err
		}
	case reflect.Int16:
		return func(str string) (reflect.Value, error) {
			v, err := strconv.ParseInt(str, 10, 16)
			return reflect.ValueOf(int16(v)), err
		}
	case reflect.Int32:
		return func(str string) (reflect.Value, error) {
			v, err := strconv.ParseInt(str, 10, 32)
			return reflect.ValueOf(int32(v)), err
		}
	case reflect.Int64:
		return func(str string) (reflect.Value, error) {
			v, err := strconv.ParseInt(str, 10, 64)
			return reflect.ValueOf(v), err
		}
	case reflect.Uint:
		return func(str string) (reflect.Value, error) {
			v, err := strconv.ParseUint(str, 10, 64)
			return reflect.ValueOf(uint(v)), err
		}
	case reflect.Uint8:
		return func(str string) (reflect.Value, error) {
			v, err := strconv.ParseUint(str, 10, 8)
			return reflect.ValueOf(uint8(v)), err
		}
	case reflect.Uint16:
		return func(str string) (reflect.Value, error) {
			v, err := strconv.ParseUint(str, 10, 16)
			return reflect.ValueOf(uint16(v)), err
		}
	case reflect.Uint32:
		return func(str string) (reflect.Value, error) {
			v, err := strconv.ParseUint(str, 10, 32)
			return reflect.ValueOf(uint32(v)), err
		}
	case reflect.Uint64:
		return func(str string) (reflect.Value, error) {
			v, err := strconv.ParseUint(str, 10, 64)
			return reflect.ValueOf(v), err
		}
	case reflect.String:
		return func(str string) (reflect.Value, error) {
			return reflect.ValueOf(str), nil
		}
	default:
		panic(fmt.Errorf("unsupported kind: %v", kind))
	}
}

func (this *Archivist) SetCurrentCollection(c Collection) {
	this.mu.Lock()
	defer this.mu.Unlock()

	compatibleVersions := make(map[string]struct{})
	for _, v := range c.CompatibleVersions() {
		compatibleVersions[v] = struct{}{}
	}

	newWrapper := &CollectionWrapper{
		Collection:         c,
		compatibleVersions: compatibleVersions,
		when:               time.Now(),
	}
	this.wrapper.Store(newWrapper)
}

func (this *Archivist) collectionWrapper() *CollectionWrapper {
	x := this.wrapper.Load()
	if x == nil {
		return nil
	}
	return x.(*CollectionWrapper)
}

func (this *Archivist) FindCollection(version string) (Collection, error) {
	wrapper := this.collectionWrapper()
	if wrapper == nil {
		return nil, errors.New("the current collection has not been set yet")
	}
	if version == "" {
		return wrapper.Collection, nil
	}
	if _, ok := wrapper.compatibleVersions[version]; ok {
		return wrapper.Collection, nil
	}

	return wrapper.Collection, ErrUpgradeNeeded
}

type Option func(arch *Archivist)

func WithLogger(log slog.Logger) Option {
	return func(arch *Archivist) {
		arch.Logger = log
	}
}

func WithRoot(root string) Option {
	return func(arch *Archivist) {
		arch.rootDir = root
	}
}

func WithWhitelist(whitelist []string) Option {
	return func(arch *Archivist) {
		arch.whitelist = make([]string, len(whitelist))
		copy(arch.whitelist, whitelist)
	}
}

func WithBlacklist(blacklist []string) Option {
	return func(arch *Archivist) {
		arch.blacklist = make([]string, len(blacklist))
		copy(arch.blacklist, blacklist)
	}
}

func WithReloadCallback(f func(newC, oldC Collection) error) Option {
	return func(arch *Archivist) {
		arch.reloadCallback = f
	}
}
