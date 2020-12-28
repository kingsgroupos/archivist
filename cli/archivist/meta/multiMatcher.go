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

package meta

import (
	"path/filepath"
	"strings"

	"github.com/kingsgroupos/archivist/cli/archivist/guesser"
	"github.com/kingsgroupos/misc"
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
