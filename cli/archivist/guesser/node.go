package guesser

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/edwingeng/deque"
	"github.com/kingsgroupos/misc"
	"github.com/kingsgroupos/misc/werr"
	"github.com/mohae/deepcopy"
	"github.com/pkg/errors"
	"github.com/xlab/treeprint"
	"go.uber.org/atomic"
)

var (
	ErrSkip         = fmt.Errorf("skip")
	ErrIncompatible = fmt.Errorf("incompatible object type detected")
)

var (
	Epsilon = .00001
)

type ValueKind int

const (
	ValueKind_Invalid ValueKind = iota
	ValueKind_Primitive
	ValueKind_Struct
	ValueKind_Map
	ValueKind_Array
	ValueKind_Ref
	ValueKind_Undetermined
)

func (vk ValueKind) String() string {
	switch vk {
	default:
		fallthrough
	case ValueKind_Invalid:
		return "invalid"
	case ValueKind_Primitive:
		return "primitive"
	case ValueKind_Struct:
		return "struct"
	case ValueKind_Map:
		return "map"
	case ValueKind_Array:
		return "array"
	case ValueKind_Ref:
		return "ref"
	case ValueKind_Undetermined:
		return "undetermined"
	}
}

type Node struct {
	Name   string
	Parent *Node

	ValueKind ValueKind
	Value     struct {
		Primitive    string
		StructFields []*Node
		MapKey       string
		MapValue     *Node
		ArrayValue   *Node
		RawRef       string
		Ref          string
	}

	Notes string
}

func NewNode(name string, parent *Node) *Node {
	return &Node{
		Name:   name,
		Parent: parent,
	}
}

func (this *Node) Tree() string {
	root := treeprint.New()
	type item struct {
		tree treeprint.Tree
		node *Node
	}

	q := deque.NewDeque()
	q.Enqueue(&item{tree: root, node: this})
	for q.Len() > 0 {
		v := q.Dequeue().(*item)
		switch v.node.ValueKind {
		case ValueKind_Invalid:
			if v.node.Name != "" {
				v.tree.AddNode(fmt.Sprintf("%s: invalid", v.node.Name))
			} else {
				v.tree.AddNode("invalid")
			}
		case ValueKind_Primitive:
			if v.node.Name != "" {
				v.tree.AddNode(fmt.Sprintf("%s: %s", v.node.Name, v.node.Value.Primitive))
			} else {
				v.tree.AddNode(v.node.Value.Primitive)
			}
		case ValueKind_Struct:
			var b1 treeprint.Tree
			if v.node.Name != "" {
				b1 = v.tree.AddBranch(fmt.Sprintf("%s: {}", v.node.Name))
			} else {
				b1 = v.tree.AddBranch("{}")
			}
			for _, field := range v.node.Value.StructFields {
				q.Enqueue(&item{tree: b1, node: field})
			}
		case ValueKind_Map:
			var b1 treeprint.Tree
			if v.node.Name != "" {
				b1 = v.tree.AddBranch(fmt.Sprintf("%s: map[%s]", v.node.Name, v.node.Value.MapKey))
			} else {
				b1 = v.tree.AddBranch(fmt.Sprintf("map[%s]", v.node.Value.MapKey))
			}
			q.Enqueue(&item{tree: b1, node: v.node.Value.MapValue})
		case ValueKind_Array:
			var b1 treeprint.Tree
			if v.node.Name != "" {
				b1 = v.tree.AddBranch(fmt.Sprintf("%s: []", v.node.Name))
			} else {
				b1 = v.tree.AddBranch("[]")
			}
			q.Enqueue(&item{tree: b1, node: v.node.Value.ArrayValue})
		case ValueKind_Ref:
			if v.node.Name != "" {
				v.tree.AddNode(fmt.Sprintf("%s: ref@%s", v.node.Name, v.node.Value.Ref))
			} else {
				v.tree.AddNode("ref@" + v.node.Value.Ref)
			}
		case ValueKind_Undetermined:
			if v.node.Name != "" {
				v.tree.AddNode(fmt.Sprintf("%s: undetermined", v.node.Name))
			} else {
				v.tree.AddNode("undetermined")
			}
		default:
			panic("impossible")
		}
	}

	lastNode := root.FindLastNode()
	if lastNode != nil {
		return lastNode.String()
	}
	return ""
}

func (this *Node) Type() string {
	switch this.ValueKind {
	case ValueKind_Primitive:
		return this.Value.Primitive
	case ValueKind_Struct:
		return "{}"
	case ValueKind_Map:
		return fmt.Sprintf("map[%s]", this.Value.MapKey)
	case ValueKind_Array:
		return "[]"
	case ValueKind_Ref:
		return "ref@" + this.Value.Ref
	case ValueKind_Undetermined:
		return "undetermined"
	default:
		panic("impossible")
	}
}

func (this *Node) Path() string {
	var sb strings.Builder
	if this.Parent != nil {
		this.Parent.pathImpl(&sb)
	} else {
		sb.WriteString("/")
	}
	if this.Name != "" {
		sb.WriteString("/")
		sb.WriteString(this.Name)
	}
	return sb.String()
}

func (this *Node) pathImpl(sb *strings.Builder) {
	if this.Parent != nil {
		this.Parent.pathImpl(sb)
	}
	if this.Name != "" {
		sb.WriteString("/")
		sb.WriteString(this.Name)
	}
	switch this.ValueKind {
	case ValueKind_Primitive:
		panic("impossible")
	case ValueKind_Struct:
	case ValueKind_Map:
		sb.WriteString("/map[]")
	case ValueKind_Array:
		sb.WriteString("/[]")
	case ValueKind_Ref:
		panic("impossible")
	case ValueKind_Undetermined:
		panic("impossible")
	default:
		panic("impossible")
	}
}

func (this *Node) AllPathTypes() []string {
	var a []string
	q := deque.NewDeque()
	q.Enqueue(this)
	for q.Len() > 0 {
		node := q.Dequeue().(*Node)
		a = append(a, node.Path(), node.Type())
		switch node.ValueKind {
		case ValueKind_Primitive:
		case ValueKind_Struct:
			for _, field := range node.Value.StructFields {
				q.Enqueue(field)
			}
		case ValueKind_Map:
			q.Enqueue(node.Value.MapValue)
		case ValueKind_Array:
			q.Enqueue(node.Value.ArrayValue)
		case ValueKind_Ref:
		case ValueKind_Undetermined:
		default:
			panic("impossible")
		}
	}

	for i := 0; i < len(a); i += 2 {
		for j := i + 2; j < len(a); j += 2 {
			if a[i] > a[j] {
				a[i], a[j] = a[j], a[i]
				a[i+1], a[j+1] = a[j+1], a[i+1]
			}
		}
	}

	return a
}

func (this *Node) parseValue(val interface{}, stringKeyTester func(nodePath string) bool) error {
	switch v := val.(type) {
	case bool:
		this.ValueKind = ValueKind_Primitive
		this.Value.Primitive = "bool"
		return nil
	case string:
		this.ValueKind = ValueKind_Primitive
		this.Value.Primitive = "string"
		return nil
	case float64:
		this.ValueKind = ValueKind_Primitive
		if math.Abs(v-math.Round(v)) >= Epsilon {
			this.Value.Primitive = "float64"
		} else {
			this.Value.Primitive = "int64"
		}
		return nil
	case map[string]interface{}:
		if len(v) == 0 {
			return werr.NewRichError(ErrSkip)
		}
		keyType := detectKeyType(v)
		if keyType == "" && stringKeyTester != nil {
			if stringKeyTester(this.Path()) {
				keyType = "string"
			}
		}
		if keyType != "" {
			one, err := mergeChildren(v)
			if err != nil {
				return errorIncompatibleWithExtraInfo(err, v)
			}
			return this.parseMap(keyType, one, stringKeyTester)
		} else {
			return this.parseStruct(v, stringKeyTester)
		}
	case []interface{}:
		if len(v) == 0 {
			this.parseUndeterminedArray()
			return nil
		}
		one, err := mergeChildren(v)
		if err != nil {
			return errorIncompatibleWithExtraInfo(err, v)
		}
		return this.parseArray(one, stringKeyTester)
	case nil:
		return werr.NewRichError(ErrSkip)
	default:
		return errors.Errorf("unexpected data type: %v", v)
	}
}

func errorIncompatibleWithExtraInfo(err *werr.RichError, v interface{}) *werr.RichError {
	if len(err.Details) < 6 {
		str := misc.ToPrettyJSONString(v)
		const maxLen = 8192
		if len(str) <= maxLen {
			return werr.NewRichError(err, oddValue(), str)
		} else {
			return werr.NewRichError(err, oddValue(), str[:maxLen-3]+"...")
		}
	} else {
		return err
	}
}

func detectKeyType(m map[string]interface{}) string {
	if len(m) == 0 {
		return ""
	}
	for k := range m {
		_, err := strconv.ParseInt(k, 0, 64)
		if err != nil {
			return ""
		}
	}
	return "int64"
}

func (this *Node) parseMap(keyType string, v interface{}, stringKeyTester func(nodePath string) bool) error {
	this.ValueKind = ValueKind_Map
	this.Value.MapKey = keyType
	newNode := NewNode("", this)
	if err := newNode.parseValue(v, stringKeyTester); err != nil {
		return err
	}
	this.Value.MapValue = newNode
	return nil
}

func mergeChildren(elements interface{}) (interface{}, *werr.RichError) {
	var err *werr.RichError
	switch es := elements.(type) {
	case map[string]interface{}:
		if len(es) == 0 {
			return nil, werr.NewRichError(ErrSkip)
		}
		var r interface{}
		for _, v := range es {
			clone := deepcopy.Copy(r)
			r, err = mergeObjects(r, v)
			if err != nil {
				if errors.Is(err, ErrSkip) {
					r = clone
					continue
				}
				return nil, errorIncompatibleWithExtraInfo(err, v)
			}
		}
		return r, nil
	case []interface{}:
		n := len(es)
		if n == 0 {
			return nil, werr.NewRichError(ErrSkip)
		}
		var r interface{}
		for i := 0; i < n; i++ {
			v := es[i]
			clone := deepcopy.Copy(r)
			r, err = mergeObjects(r, v)
			if err != nil {
				if errors.Is(err, ErrSkip) {
					r = clone
					continue
				}
				return nil, errorIncompatibleWithExtraInfo(err, v)
			}
		}
		return r, nil
	default:
		panic("impossible")
	}
}

var (
	odd atomic.Int64
)

func oddValue() string {
	return "odd" + fmt.Sprint(odd.Inc())
}

func mergeObjects(obj1, obj2 interface{}) (interface{}, *werr.RichError) {
	if obj1 == nil {
		clone := deepcopy.Copy(obj2)
		return clone, nil
	}
	if obj2 == nil {
		clone := deepcopy.Copy(obj1)
		return clone, nil
	}
	if reflect.TypeOf(obj1) != reflect.TypeOf(obj2) {
		return nil, werr.NewRichError(ErrIncompatible, oddValue(),
			fmt.Errorf("obj1: %s, obj2: %s", misc.ToPrettyJSONString(obj1), misc.ToPrettyJSONString(obj2)))
	}

	switch o1 := obj1.(type) {
	case map[string]interface{}:
		if o2, ok := obj2.(map[string]interface{}); ok {
			for k, v2 := range o2 {
				if v1, ok := o1[k]; ok {
					var err *werr.RichError
					clone := deepcopy.Copy(v1)
					o1[k], err = mergeObjects(v1, v2)
					if err != nil {
						if errors.Is(err, ErrSkip) {
							o1[k] = clone
							continue
						}
						return nil, err
					}
				} else {
					o1[k] = v2
				}
			}
		}
		return o1, nil
	case []interface{}:
		if o2, ok := obj2.([]interface{}); ok {
			one, err := mergeChildren(append(o1, o2...))
			if err != nil {
				return nil, err
			}
			return []interface{}{one}, nil
		} else {
			one, err := mergeChildren(o1)
			if err != nil {
				return nil, err
			}
			return []interface{}{one}, nil
		}
	case float64:
		f1 := obj1.(float64)
		if math.Abs(f1-math.Round(f1)) >= Epsilon {
			return obj1, nil
		}
		f2 := obj2.(float64)
		if math.Abs(f2-math.Round(f2)) >= Epsilon {
			return obj2, nil
		}
	}

	clone := deepcopy.Copy(obj1)
	return clone, nil
}

func (this *Node) parseUndeterminedArray() {
	this.ValueKind = ValueKind_Array
	newNode := NewNode("", this)
	newNode.ValueKind = ValueKind_Undetermined
	this.Value.ArrayValue = newNode
}

func (this *Node) parseArray(v interface{}, stringKeyTester func(nodePath string) bool) error {
	this.ValueKind = ValueKind_Array
	newNode := NewNode("", this)
	if err := newNode.parseValue(v, stringKeyTester); err != nil {
		return err
	}
	this.Value.ArrayValue = newNode
	return nil
}

func (this *Node) parseStruct(m map[string]interface{}, stringKeyTester func(nodePath string) bool) error {
	this.ValueKind = ValueKind_Struct
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if _, ok := m["ID"]; ok {
		idx := sort.SearchStrings(keys, "ID")
		copy(keys[1:], keys[:idx])
		keys[0] = "ID"
	}

	for _, k := range keys {
		newNode := NewNode(k, this)
		if err := newNode.parseValue(m[k], stringKeyTester); err != nil {
			if errors.Is(err, ErrSkip) {
				continue
			}
			return err
		}
		this.Value.StructFields = append(this.Value.StructFields, newNode)
	}
	if len(this.Value.StructFields) == 0 {
		return werr.NewRichError(ErrSkip)
	}
	return nil
}

func PrettyErrorIncompatible(err error) string {
	var sb strings.Builder
	var number int
	prettyErrorIncompatibleImpl(err, &number, &sb)
	return strings.TrimSuffix(sb.String(), "\n")
}

func prettyErrorIncompatibleImpl(err error, number *int, sb *strings.Builder) {
	split := func() {
		*number++
		_, _ = fmt.Fprintf(sb, "%s [%02d] =\n", strings.Repeat("=", 30), *number)
	}
	if richErr, ok := err.(*werr.RichError); ok {
		prettyErrorIncompatibleImpl(richErr.Unwrap(), number, sb)
		for i := 0; i < len(richErr.Details); i += 2 {
			split()
			_, _ = fmt.Fprintf(sb, "%s\n", richErr.Details[i+1])
		}
	} else {
		split()
		_, _ = fmt.Fprintf(sb, "%+v\n", err)
	}
}

func (this *Node) Traverse(f func(node *Node) bool) {
	this.traverseImpl(f)
}

func (this *Node) traverseImpl(f func(node *Node) bool) {
	if !f(this) {
		return
	}
	switch this.ValueKind {
	case ValueKind_Primitive:
	case ValueKind_Struct:
		for _, field := range this.Value.StructFields {
			field.traverseImpl(f)
		}
	case ValueKind_Map:
		this.Value.MapValue.traverseImpl(f)
	case ValueKind_Array:
		this.Value.ArrayValue.traverseImpl(f)
	case ValueKind_Ref:
	case ValueKind_Undetermined:
	default:
		panic("impossible")
	}
}

func (this *Node) Meaningless() bool {
	switch this.ValueKind {
	case ValueKind_Primitive:
		return false
	case ValueKind_Struct:
		for _, field := range this.Value.StructFields {
			if field.Meaningless() {
				continue
			}
			return false
		}
		return true
	case ValueKind_Map:
		return this.Value.MapValue.Meaningless()
	case ValueKind_Array:
		return this.Value.ArrayValue.Meaningless()
	case ValueKind_Ref:
		return false
	case ValueKind_Undetermined:
		return true
	default:
		panic("impossible")
	}
}

func (this *Node) RemoveUndeterminedChildren() {
	switch this.ValueKind {
	case ValueKind_Primitive:
	case ValueKind_Struct:
		newFields := make([]*Node, 0, len(this.Value.StructFields))
		for _, field := range this.Value.StructFields {
			if field.Meaningless() {
				continue
			}
			newFields = append(newFields, field)
		}
		if len(newFields) != cap(newFields) {
			this.Value.StructFields = newFields
		}
	case ValueKind_Map:
		this.Value.MapValue.RemoveUndeterminedChildren()
	case ValueKind_Array:
		this.Value.ArrayValue.RemoveUndeterminedChildren()
	case ValueKind_Ref:
	case ValueKind_Undetermined:
	default:
		panic("impossible")
	}
}
