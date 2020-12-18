package meta

import (
	"path/filepath"
	"strings"

	"gitlab-ee.funplus.io/watcher/watcher/archivist/cli/archivist/guesser"
	"gitlab-ee.funplus.io/watcher/watcher/misc"
)

type MultiMatcher struct {
	f1, f2 string
	p1, p2 *Parser
}

func NewMultiMatcher() *MultiMatcher {
	return &MultiMatcher{}
}

func (this *MultiMatcher) ParseMetaFiles(jsonFile string) error {
	dir1 := filepath.Dir(jsonFile)
	dir2 := filepath.Dir(dir1)
	dir3 := filepath.Dir(dir2)
	if strings.HasSuffix(dir1, "/.runtime") {
		jsonFile = filepath.Join(dir2, filepath.Base(jsonFile))
	} else if strings.HasSuffix(dir2, "/.runtime") {
		jsonFile = filepath.Join(dir3, filepath.Base(jsonFile))
	}

	f0 := strings.TrimSuffix(jsonFile, ".json")
	f1 := f0 + ".meta"
	f2 := f0 + ".suggested.meta"

	if err := misc.FindFile(f1); err == nil {
		this.p1 = NewParser()
		if err := this.p1.ParseFile(f1); err != nil {
			return err
		}
	}
	if err := misc.FindFile(f2); err == nil {
		this.p2 = NewParser()
		if err := this.p2.ParseFile(f2); err != nil {
			return err
		}
	}

	this.f1, this.f2 = f1, f2
	return nil
}

func (this *MultiMatcher) Match(node *guesser.Node) (*Meta, bool) {
	nodePath := node.Path()
	if this.p1 != nil {
		mt, ok := this.p1.Match(nodePath)
		if ok {
			return mt, true
		}
	}
	if this.p2 != nil {
		mt, ok := this.p2.Match(nodePath)
		if ok {
			return mt, true
		}
	}
	if this.p1 != nil {
		mt, ok := this.p1.WildMatch(nodePath)
		if ok {
			return mt, true
		}
	}
	if this.p2 != nil {
		mt, ok := this.p2.WildMatch(nodePath)
		if ok {
			return mt, true
		}
	}

	return nil, false
}

func (this *MultiMatcher) StringKeyTester() func(nodePath string) bool {
	m := make(map[string]*Meta)
	a := make([]*Meta, 0)
	if this.p1 != nil {
		for _, mt := range this.p1.metaMap {
			if mt.NewType == "map[string]" {
				m[mt.Path] = mt
				a = append(a, mt)
			}
		}
	}
	if this.p2 != nil {
		for _, mt := range this.p2.metaMap {
			if mt.NewType == "map[string]" {
				if _, ok := m[mt.Path]; ok {
					continue
				}
				m[mt.Path] = mt
				a = append(a, mt)
			}
		}
	}
	return func(nodePath string) bool {
		if mt, ok := m[nodePath]; ok {
			mt.Counter++
			return true
		}
		_, matched := wildMatch(a, nodePath)
		return matched
	}
}

func (this *MultiMatcher) UnusedMetas() map[string][]*Meta {
	m := make(map[string][]*Meta)
	if this.p1 != nil {
		a := this.p1.UnusedMetas()
		if len(a) > 0 {
			m[this.f1] = a
		}
	}
	if this.p2 != nil {
		a := this.p2.UnusedMetas()
		if len(a) > 0 {
			m[this.f2] = a
		}
	}
	return m
}
