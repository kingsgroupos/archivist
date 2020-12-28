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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kingsgroupos/misc"
	"github.com/pkg/errors"
)

const EnvName_ConfGroup = "WATCHER_CONF_GROUP"
const EnvName_ConfSubgroup = "WATCHER_CONF_SUBGROUP"

var (
	errEnv = fmt.Errorf("environmental variable %s should not be empty", EnvName_ConfGroup)
)

func ConfGroup() string {
	return os.Getenv(EnvName_ConfGroup)
}

func ConfSubgroup() string {
	return os.Getenv(EnvName_ConfSubgroup)
}

func MustHaveConfGroup() {
	if ConfGroup() == "" {
		panic(errEnv)
	}
}

func Dump(c Collection, outputDir string) error {
	if c == nil {
		return errors.New("c cannot be nil")
	}
	if err := misc.FindDirectory(outputDir); err != nil {
		return errors.WithStack(err)
	}

	for k, v := range c.Filename2Conf() {
		if v == nil {
			continue
		}
		data1, err := json.Marshal(v)
		if err != nil {
			return errors.Errorf("json.Marshal failed. file: %s, err: %+v", k, err)
		}
		var obj interface{}
		if err := json.Unmarshal(data1, &obj); err != nil {
			return errors.Errorf("json.Unmarshal failed. file: %s, err: %+v", k, err)
		}
		data2, err := json.MarshalIndent(obj, "", "\t")
		if err != nil {
			return errors.Errorf("json.MarshalIndent failed. file: %s, err: %+v", k, err)
		}

		data2 = append(data2, '\n')
		file := filepath.Join(outputDir, k)
		if err := ioutil.WriteFile(file, data2, 0644); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
