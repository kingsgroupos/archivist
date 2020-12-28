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
	"os"
	"path/filepath"
	"strings"

	"github.com/kingsgroupos/misc"
	"github.com/kingsgroupos/misc/variable"
	"github.com/kingsgroupos/misc/wildcard"
	"github.com/spf13/cobra"
)

const (
	easyjsonSuffix = "_easyjson.go"
)

var orphanCmd orphanCmdT

var orphanCmdCobra = &cobra.Command{
	Use:   "orphan <dataDir> <codeDir> <codeFileExt> [ignore1,ignore2,...]",
	Short: "Find (and delete, optionally) orphan files in your code",
	Run:   orphanCmd.execute,
}

func init() {
	rootCmd.AddCommand(orphanCmdCobra)
	cmd := orphanCmdCobra
	cmd.Args = cobra.RangeArgs(3, 4)
	cmd.Flags().StringVar(&orphanCmd.structNameSuffix,
		"structNameSuffix", "Conf", "name suffix of struct")
	cmd.Flags().StringSliceVar(&orphanCmd.accompanyingFileSuffixes,
		"accompanyingFileSuffix", nil, "delete accompanying files when -d is specified")
	cmd.Flags().BoolVarP(&orphanCmd.delete,
		"delete", "d", false, "delete orphan files")
}

type orphanCmdT struct {
	structNameSuffix         string
	accompanyingFileSuffixes []string
	delete                   bool
}

func (this *orphanCmdT) execute(cmd *cobra.Command, args []string) {
	for i := 0; i < 2; i++ {
		if err := misc.FindDirectory(args[i]); err != nil {
			panic(err)
		}
	}

	dataDir := args[0]
	codeDir := args[1]
	codeExt := "." + strings.TrimLeft(args[2], ".")
	ignores := []string{TplCollection, TplCollectionExtension}
	if len(args) > 3 {
		sep := "|"
		if !strings.Contains(args[3], "|") && strings.Contains(args[3], ",") {
			sep = ","
		}
		a := misc.Split(args[3], sep)
		ignores = append(ignores, a...)
	}

	if codeExt == ".go" {
		if misc.IndexStrings(this.accompanyingFileSuffixes, easyjsonSuffix) < 0 {
			this.accompanyingFileSuffixes = append(this.accompanyingFileSuffixes, easyjsonSuffix)
		}
	}

	pats1 := []string{
		filepath.Join(dataDir, "*.json"),
		filepath.Join(dataDir, "*.js"),
	}
	allDataFiles, err := misc.AllFiles(pats1, "", true)
	if err != nil {
		panic(err)
	}

	pats2 := []string{
		filepath.Join(codeDir, "*"+codeExt),
	}
	allCodeFiles, err := misc.AllFiles(pats2, "", true)
	if err != nil {
		panic(err)
	}

	m := make(map[string]struct{})
	for _, file := range allDataFiles {
		basename1 := filepath.Base(file)
		var basename2 string
		switch {
		case strings.HasSuffix(basename1, ".json"):
			basename2 = strings.TrimSuffix(basename1, ".json")
		case strings.HasSuffix(basename1, ".js"):
			basename2 = strings.TrimSuffix(basename1, ".js")
		default:
			panic(fmt.Errorf("unexpected file name: %s", basename1))
		}
		m[variable.ToCamel(basename2)] = struct{}{}
	}

	var n int
outer:
	for _, file := range allCodeFiles {
		for _, suffix := range this.accompanyingFileSuffixes {
			if strings.HasSuffix(file, suffix) {
				continue outer
			}
		}

		basename1 := filepath.Base(file)
		basename2 := strings.TrimSuffix(basename1, codeExt)
		basename3 := strings.TrimSuffix(basename2, this.structNameSuffix)
		if _, ok := m[basename3]; ok {
			continue
		}
		for _, str := range ignores {
			if str == basename2 || str == basename2+codeExt || wildcard.Match(str, basename2) {
				continue outer
			}
		}

		if !this.delete {
			_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("orphan: "+file))
			n++
		} else {
			this.deleteOrphanFile(file)
		}

		for _, suffix := range this.accompanyingFileSuffixes {
			accompanyingFile := filepath.Join(codeDir, basename2+suffix)
			if !this.delete {
				this.hintOrphanFile(accompanyingFile)
			} else {
				this.deleteOrphanFile(accompanyingFile)
			}
		}
	}
	if n > 0 {
		os.Exit(1)
	}
}

func (this *orphanCmdT) hintOrphanFile(file string) {
	if misc.FindFile(file) == nil {
		_, _ = fmt.Fprintf(os.Stderr, "orphan: %s\n", file)
	}
}

func (this *orphanCmdT) deleteOrphanFile(file string) {
	if misc.FindFile(file) != nil {
		return
	}
	if err := os.Remove(file); err != nil {
		panic(err)
	}
	_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("orphan deleted: "+file))
}
