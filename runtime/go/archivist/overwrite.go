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

package archivist

import (
	"sort"

	"github.com/pkg/errors"
)

type Overwrite struct {
	FileLevel bool
	Target    string
	Data      []byte
}

type organizedOverwrites struct {
	fileLevel    map[string][][]byte
	contentLevel map[string][][]byte
}

func (this *organizedOverwrites) organize(c Collection, o []Overwrite) error {
	if len(o) == 0 {
		return nil
	}

	m1 := make(map[string][][]byte)
	m2 := make(map[string][][]byte)
	targets := c.Filename2Conf()
	for _, x := range o {
		if _, ok := targets[x.Target]; !ok {
			return errors.Errorf("<archivist> unknown overwrite target: %s", x.Target)
		}
		if x.FileLevel {
			m1[x.Target] = append(m1[x.Target], x.Data)
		} else {
			m2[x.Target] = append(m2[x.Target], x.Data)
		}
	}

	this.fileLevel = m1
	this.contentLevel = m2
	return nil
}

func (this *organizedOverwrites) findAffected(revRefGraph map[string][]string) []string {
	var pending []string
	for f := range this.fileLevel {
		pending = append(pending, f)
	}
	for f := range this.contentLevel {
		pending = append(pending, f)
	}

	m := make(map[string]struct{})
	affected := make(map[string]struct{})
	for n := len(pending); n > 0; n = len(pending) {
		f := pending[n-1]
		pending = pending[:n-1]
		if _, ok := m[f]; ok {
			continue
		}

		m[f] = struct{}{}
		for _, referer := range revRefGraph[f] {
			affected[referer] = struct{}{}
			pending = append(pending, referer)
		}
	}

	for f := range this.fileLevel {
		delete(affected, f)
	}
	for f := range this.contentLevel {
		delete(affected, f)
	}

	a := make([]string, 0, len(affected))
	for f := range affected {
		a = append(a, f)
	}
	sort.Strings(a)
	return a
}
