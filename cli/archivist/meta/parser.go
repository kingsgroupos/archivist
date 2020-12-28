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
	"bufio"
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/kingsgroupos/misc"
	"github.com/kingsgroupos/misc/werr"
	"github.com/kingsgroupos/misc/wildcard"
	"github.com/pkg/errors"
)

type Parser struct {
	metaMap   map[string]*Meta
	wildMetas []*Meta
}

func NewParser() *Parser {
	return &Parser{
		metaMap: make(map[string]*Meta),
	}
}

func (this *Parser) ParseFile(file string) error {
	bts, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = this.ParseBytes(bts)
	if err != nil {
		return werr.NewRichError(err, "file", file)
	}
	return nil
}

func (this *Parser) ParseBytes(bts []byte) error {
	var lineNumber int
	scanner := bufio.NewScanner(bytes.NewReader(bts))
	for scanner.Scan() {
		lineNumber++
		meta, err := parseLine(scanner.Bytes())
		if err != nil {
			return werr.NewRichError(err, "line", lineNumber)
		}
		if meta == nil {
			continue
		}
		if _, ok := this.metaMap[meta.Path]; ok {
			return werr.NewRichError(errors.New("path duplication detected"), "line", lineNumber)
		}
		this.metaMap[meta.Path] = meta
		if strings.ContainsAny(meta.Path, "*?") {
			this.wildMetas = append(this.wildMetas, meta)
		}
	}

	return nil
}

func wildMatch(metas []*Meta, nodePath string) (*Meta, bool) {
	fields := misc.Split(nodePath, "/")
	numFields := len(fields)

outer:
	for _, mt := range metas {
		if len(mt.PathFields) != numFields {
			continue
		}
		for i := 0; i < numFields; i++ {
			if !wildcard.Match(mt.PathFields[i], fields[i]) {
				continue outer
			}
		}
		mt.Counter++
		return mt, true
	}

	return nil, false
}

func (this *Parser) Match(nodePath string) (*Meta, bool) {
	if mt, ok := this.metaMap[nodePath]; ok {
		mt.Counter++
		return mt, true
	}
	return nil, false
}

func (this *Parser) WildMatch(nodePath string) (*Meta, bool) {
	return wildMatch(this.wildMetas, nodePath)
}

func (this *Parser) UnusedMetas() []*Meta {
	var a []*Meta
	for _, mt := range this.metaMap {
		if mt.Counter == 0 {
			a = append(a, mt)
		}
	}
	return a
}
