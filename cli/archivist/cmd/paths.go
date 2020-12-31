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
	"os"
	"strings"
	"text/tabwriter"

	"github.com/kingsgroupos/archivist/cli/archivist/guesser"
	"github.com/spf13/cobra"
)

var pathsCmd pathsCmdT

var pathsCmdCobra = &cobra.Command{
	Use:   "paths <file>",
	Short: "Show all the node paths and data types of a .json/.js file",
	Run:   pathsCmd.execute,
}

func init() {
	rootCmd.AddCommand(pathsCmdCobra)
	cmd := pathsCmdCobra
	cmd.Flags().BoolVar(&pathsCmd.bare,
		"bare", false, "ignore meta files")

	pathsCmd.registerSharedFlags(cmd)
}

type pathsCmdT struct {
	bare bool

	sharedFlags
}

func (this *pathsCmdT) execute(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		panic("no input file")
	}

	this.validateShared()

	file := args[0]
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	switch {
	case strings.HasSuffix(file, ".json"):
	case strings.HasSuffix(file, ".js"):
		data = evalJavascript(data, file)
	default:
		panic("the input file must be a .json or a .js")
	}

	var g *guesser.Guesser
	if this.bare {
		g = guesser.NewGuesser()
		if err := g.ParseBytes(data); err != nil {
			panic(guesser.PrettyErrorIncompatible(err))
		}
	} else if strings.HasSuffix(file, ".json") {
		g = this.buildGuesser(file, nil, "")
	} else {
		g = this.buildGuesserWithJavascriptFile(data, file, nil, "")
	}

	printAllPathTypes(g.Root.AllPathTypes())
}

func printAllPathTypes(all []string) {
	w := tabwriter.NewWriter(os.Stdout, 10, 4, 3, ' ', 0)
	_, _ = fmt.Fprintln(w, "Path\tType")
	for i := 0; i < len(all); i += 2 {
		_, _ = fmt.Fprintf(w, "%s\t%s\n", all[i], all[i+1])
	}
	_ = w.Flush()
}
