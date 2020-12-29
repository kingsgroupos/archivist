// BSD 3-Clause License
//
// Copyright (c) 2020, Kingsgroup
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kingsgroupos/archivist/cli/archivist/guesser"
	"github.com/kingsgroupos/archivist/cli/archivist/meta"
	"github.com/kingsgroupos/misc"
	"github.com/kingsgroupos/misc/variable"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type sharedFlags struct {
	intType        string
	floatType      string
	floatPrecision int

	skipped map[string]struct{}
}

func (this *sharedFlags) registerSharedFlags(cmd *cobra.Command) {
	const allowedIntTypes = "int, int8, int16, int32, int64, uint, uint8, uint16 or uint32"
	cmd.Flags().StringVar(&this.intType,
		"intType", "int64", "the default integer type, must be "+allowedIntTypes)
	cmd.Flags().StringVar(&this.floatType,
		"floatType", "float64", "the default floating number type, must be float32 or float64")
	cmd.Flags().IntVar(&this.floatPrecision,
		"floatPrecision", 5, "the precision of floating numbers")
}

func (this *sharedFlags) validateShared() {
	validateIntType(this.intType)
	validateFloatType(this.floatType)
	validateFloatPrecision(this.floatPrecision)
}

func validateIntType(intType string) {
	switch intType {
	case "int":
	case "int8":
	case "int16":
	case "int32":
	case "int64":
	case "uint":
	case "uint8":
	case "uint16":
	case "uint32":
	default:
		panic(fmt.Sprintf("invalid --intType: %s", intType))
	}
}

func validateFloatType(floatType string) {
	switch floatType {
	case "float32":
	case "float64":
	default:
		panic(fmt.Sprintf("invalid --floatType: %s", floatType))
	}
}

func validateFloatPrecision(n int) {
	if n <= 0 || n > 8 {
		panic("--floatPrecision must be in between [1, 8]")
	}

	guesser.Epsilon = math.Pow(10, -float64(n))
}

func (this *sharedFlags) buildGuesser(jsonFile string, primaryStructNameMap map[string]int, structNameSuffix string) *guesser.Guesser {
	var err error
	matcher := meta.NewMultiMatcher()
	err = matcher.ParseMetaFiles(jsonFile)
	if err != nil {
		panic(err)
	}

	stringKeyTester := matcher.StringKeyTester()
	g := guesser.NewGuesser(guesser.WithStringKeyTester(stringKeyTester))
	err = g.ParseFile(jsonFile)
	if err != nil {
		panic(guesser.PrettyErrorIncompatible(err))
	}

	return this.buildGuesserImpl(g, matcher, primaryStructNameMap, structNameSuffix)
}

func (this *sharedFlags) buildGuesserWithJavascriptFile(data []byte, file string, primaryStructNameMap map[string]int, structNameSuffix string) *guesser.Guesser {
	var err error
	matcher := meta.NewMultiMatcher()
	jsonFile := strings.TrimSuffix(file, ".js") + ".json"
	err = matcher.ParseMetaFiles(jsonFile)
	if err != nil {
		panic(err)
	}

	stringKeyTester := matcher.StringKeyTester()
	g := guesser.NewGuesser(guesser.WithStringKeyTester(stringKeyTester))
	err = g.ParseBytes(data)
	if err != nil {
		panic(guesser.PrettyErrorIncompatible(err))
	}

	return this.buildGuesserImpl(g, matcher, primaryStructNameMap, structNameSuffix)
}

func (this *sharedFlags) buildGuesserImpl(g *guesser.Guesser, matcher *meta.MultiMatcher, primaryStructNameMap map[string]int, structNameSuffix string) *guesser.Guesser {
	g.Root.Traverse(func(node *guesser.Node) bool {
		if mt, ok := matcher.Match(node); ok {
			node.Notes = mt.Notes
			err := this.changeType(node, mt, primaryStructNameMap, structNameSuffix)
			if err != nil {
				panic(err)
			}
		} else {
			if node.ValueKind == guesser.ValueKind_Primitive {
				switch node.Value.Primitive {
				case "int64":
					node.Value.Primitive = this.intType
				case "float64":
					node.Value.Primitive = this.floatType
				}
			}
		}
		return true
	})

	unused := matcher.UnusedMetas()
	if len(unused) > 0 {
		var sb strings.Builder
		sb.WriteString("WARNING: invalid meta config detected.\n")
		for f, a := range unused {
			sb.WriteString("==============================\n")
			_, _ = fmt.Fprintf(&sb, "%s\n", f)
			sb.WriteString("==============================\n")
			for _, mt := range a {
				_, _ = fmt.Fprintf(&sb, "%s\n", mt.RawData)
			}
		}
		fmt.Println(sb.String())
	}

	return g
}

func (this *sharedFlags) changeType(node *guesser.Node, mt *meta.Meta, primaryStructNameMap map[string]int, structNameSuffix string) error {
	const errFmt = "incompatible type detected. path: %s, actualType: %s, newType: %s"
	switch node.ValueKind {
	case guesser.ValueKind_Primitive:
		nodeType := node.Type()
		switch nodeType {
		case "bool":
			switch mt.NewType {
			case "bool":
			default:
				return errors.Errorf(errFmt, mt.Path, nodeType, mt.NewType)
			}
		case "int64":
			switch mt.NewType {
			case "int":
			case "int8":
			case "int16":
			case "int32":
			case "int64":
			case "uint":
			case "uint8":
			case "uint16":
			case "uint32":
			case "float32":
			case "float64":
			default:
				if strings.HasPrefix(mt.NewType, "ref@") {
					rawRef := strings.TrimPrefix(mt.NewType, "ref@")
					ref := rawRef + structNameSuffix
					if _, ok := this.skipped[rawRef]; ok {
						mt.NewType = "int64"
						break
					}
					camelRef := variable.ToCamel(ref)
					if primaryStructNameMap != nil {
						if _, ok := primaryStructNameMap[camelRef]; !ok {
							return errors.Errorf("unknown reference: %s", ref)
						}
						primaryStructNameMap[camelRef]++
					}
					allowed := func() bool {
						strs := misc.Split(node.Path(), "/")
						for _, s := range strs {
							if !strings.HasSuffix(s, "]") {
								return true
							}
						}
						return false
					}
					if !allowed() {
						return errors.Errorf("reference (ref@) must reside in a struct. node: %s", node.Path())
					}
					node.ValueKind = guesser.ValueKind_Ref
					node.Value.RawRef = rawRef
					node.Value.Ref = ref
					return nil
				}
				return errors.Errorf(errFmt, mt.Path, nodeType, mt.NewType)
			}
		case "string":
			switch mt.NewType {
			case "string":
			case "datetime":
			case "duration":
			case "i18n":
			default:
				return errors.Errorf(errFmt, mt.Path, nodeType, mt.NewType)
			}
		case "float64":
			switch mt.NewType {
			case "float32":
			case "float64":
			default:
				return errors.Errorf(errFmt, mt.Path, nodeType, mt.NewType)
			}
		default:
			panic("impossible")
		}
		if nodeType != mt.NewType {
			node.Value.Primitive = mt.NewType
		}
	case guesser.ValueKind_Struct:
		nodeType := node.Type()
		switch mt.NewType {
		case "{}":
		case "map[string]":
		default:
			return errors.Errorf(errFmt, mt.Path, nodeType, mt.NewType)
		}
		if nodeType != mt.NewType {
			panic("impossible")
		}
	case guesser.ValueKind_Map:
		nodeType := node.Type()
		switch mt.NewType {
		case "map[int]":
		case "map[int8]":
		case "map[int16]":
		case "map[int32]":
		case "map[int64]":
		case "map[uint]":
		case "map[uint8]":
		case "map[uint16]":
		case "map[uint32]":
		case "map[string]":
		default:
			return errors.Errorf(errFmt, mt.Path, nodeType, mt.NewType)
		}
		if nodeType != mt.NewType {
			node.Value.MapKey = mt.NewType[4 : len(mt.NewType)-1]
		}
	case guesser.ValueKind_Array:
		nodeType := node.Type()
		switch mt.NewType {
		case "[]":
		default:
			return errors.Errorf(errFmt, mt.Path, nodeType, mt.NewType)
		}
		if nodeType != mt.NewType {
			panic("impossible")
		}
	case guesser.ValueKind_Undetermined:
		switch mt.NewType {
		case "bool":
		case "int":
		case "int8":
		case "int16":
		case "int32":
		case "int64":
		case "uint":
		case "uint8":
		case "uint16":
		case "uint32":
		case "float32":
		case "float64":
		case "string":
		case "datetime":
		case "duration":
		case "i18n":
		default:
			if strings.HasPrefix(mt.NewType, "ref@") {
				node.ValueKind = guesser.ValueKind_Primitive
				node.Value.Primitive = "int64"
				return this.changeType(node, mt, primaryStructNameMap, structNameSuffix)
			}
			return errors.Errorf(errFmt, mt.Path, node.Type(), mt.NewType)
		}
		node.ValueKind = guesser.ValueKind_Primitive
		node.Value.Primitive = mt.NewType
	default:
		panic("impossible")
	}

	return nil
}

func loadSkipped(allFiles []string) map[string]struct{} {
	skipped := make(map[string]struct{})
	uniqueDirs1 := make(map[string]struct{})
	for _, f := range allFiles {
		d := filepath.Dir(f)
		if _, ok := uniqueDirs1[d]; ok {
			continue
		}
		uniqueDirs1[d] = struct{}{}
	}

	var uniqueDirs2 []string
	for d := range uniqueDirs1 {
		uniqueDirs2 = append(uniqueDirs2, d)
	}
	sort.Strings(uniqueDirs2)

	for _, d := range uniqueDirs2 {
		a, err := loadSkippedImpl(d)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			panic(err)
		}
		for _, x := range a {
			skipped[x] = struct{}{}
		}
		break
	}

	return skipped
}

func loadSkippedImpl(dir string) ([]string, error) {
	file := filepath.Join(dir, "conf.skipped")
	if err := misc.FindFile(file); err != nil {
		return nil, err
	}
	return misc.AllFileLines(file)
}

func evalJavascript(data []byte, file string) []byte {
	if len(data) == 0 {
		panic(fmt.Errorf("%s is empty", file))
	}

	var script bytes.Buffer
	_, _ = fmt.Fprintf(&script, "var module = {}\n")
	_, _ = fmt.Fprintf(&script, "%s\n", data)
	_, _ = fmt.Fprintf(&script, "var json = JSON.stringify(module.exports)\n")
	_, _ = fmt.Fprintf(&script, "console.log(json)\n")

	f, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	filePath := f.Name()
	defer func() {
		_ = f.Close()
		_ = os.Remove(filePath)
	}()

	if _, err := f.Write(script.Bytes()); err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.Command("node", filePath)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(stderr.String())
		panic(err)
	}

	invalid := fmt.Errorf("%s is NOT a valid config file", file)
	out := bytes.Trim(stdout.Bytes(), " \t\n")
	if n := len(out); n == 0 {
		panic(invalid)
	} else if n < 2 {
		panic(invalid)
	} else if s, e := out[0], out[n-1]; (s != '{' && s != '[') || (e != '}' && e != ']') {
		panic(invalid)
	}
	return out
}
