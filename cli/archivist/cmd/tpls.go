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
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/kingsgroupos/archivist/cli/archivist/g"
	"github.com/kingsgroupos/misc"
	"github.com/spf13/cobra"
)

const (
	TplStruct              = "struct"
	TplCollection          = "collection"
	TplCollectionExtension = "collectionExtension"
)

var TplMap = map[string]string{
	TplStruct:              g.TemplateStruct,
	TplCollection:          g.TemplateCollection,
	TplCollectionExtension: g.TemplateCollectionExtension,
}

var tplsCmd tplsCmdT

var tplsCmdCobra = &cobra.Command{
	Use:   "tpls",
	Short: "Output the default code templates",
	Run:   tplsCmd.execute,
}

func init() {
	rootCmd.AddCommand(tplsCmdCobra)
	cmd := tplsCmdCobra
	cmd.Flags().StringVarP(&tplsCmd.outputDir,
		"outputDir", "o", "", "the output directory")
	cmd.Flags().BoolVarP(&tplsCmd.force,
		"force", "f", false, "overwrite existing files")
}

type tplsCmdT struct {
	outputDir string
	force     bool
}

func (this *tplsCmdT) execute(cmd *cobra.Command, args []string) {
	var a []string
	for k, v := range TplMap {
		a = append(a, k, v)
	}

	if this.outputDir == "" {
		for i := 0; i < len(a); i += 2 {
			fmt.Printf("=== %s %s\n\n%s\n\n", a[i], strings.Repeat("=", 60-len(a[i])), a[i+1])
		}
		return
	}

	if err := misc.FindDirectory(this.outputDir); err != nil {
		panic(fmt.Errorf("%s does not exist", this.outputDir))
	}

	for i := 0; i < len(a); i += 2 {
		fPath := filepath.Join(this.outputDir, a[i]+".tpl")
		if err := misc.FindFile(fPath); err == nil {
			if !this.force {
				panic(fmt.Errorf("%s.tpl already exists, add -f to overwrite it", a[i]))
			}
		}
	}

	for i := 0; i < len(a); i += 2 {
		fPath := filepath.Join(this.outputDir, a[i]+".tpl")
		trimmed := strings.TrimLeft(a[i+1], "\r\n")
		if err := ioutil.WriteFile(fPath, []byte(trimmed), 0644); err != nil {
			panic(err)
		}
		fmt.Println(fPath)
	}
}
