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
	"bytes"
	"strings"

	"github.com/kingsgroupos/misc"
	"github.com/pkg/errors"
)

type Meta struct {
	Path    string
	NewType string
	Notes   string
	RawData []byte
	Counter int

	PathFields []string
}

var (
	validTypes = map[string]struct{}{
		"bool":     {},
		"int":      {},
		"int8":     {},
		"int16":    {},
		"int32":    {},
		"int64":    {},
		"uint":     {},
		"uint8":    {},
		"uint16":   {},
		"uint32":   {},
		"float32":  {},
		"float64":  {},
		"string":   {},
		"datetime": {},
		"duration": {},
		"i18n":     {},
		"[]":       {},
		"{}":       {},
	}
)

func parseLine(line []byte) (*Meta, error) {
	trimmed := bytes.TrimSpace(line)
	if len(trimmed) == 0 {
		return nil, nil
	}
	fields := bytes.Fields(line)
	if len(fields) < 2 {
		return nil, errors.New("too few fields")
	}
	if len(fields) > 3 {
		fields[2] = bytes.Join(fields[2:], []byte(" "))
		fields = fields[:3]
	}

	mt := &Meta{
		Path:    string(fields[0]),
		NewType: string(fields[1]),
		RawData: append([]byte{}, line...),
	}
	if len(fields) > 2 {
		mt.Notes = string(fields[2])
	}
	if !isValidNewType(mt.NewType) {
		return nil, errors.Errorf("invalid new type: %s. path: %s", mt.NewType, mt.Path)
	}

	mt.PathFields = misc.Split(mt.Path, "/")
	return mt, nil
}

func isValidNewType(t string) bool {
	if _, ok := validTypes[t]; ok {
		return true
	} else if strings.HasPrefix(t, "map[") && strings.HasSuffix(t, "]") {
		str := t[4 : len(t)-1]
		if _, ok := validTypes[str]; ok {
			return true
		}
	} else if strings.HasPrefix(t, "ref@") {
		return true
	}
	return false
}
