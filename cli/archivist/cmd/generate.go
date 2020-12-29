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
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"text/template"

	"github.com/edwingeng/deque"
	"github.com/kingsgroupos/archivist/cli/archivist/guesser"
	"github.com/kingsgroupos/misc"
	"github.com/kingsgroupos/misc/chksum"
	"github.com/kingsgroupos/misc/variable"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var generateCmd generateCmdT

var generateCmdCobra = &cobra.Command{
	Use:   "generate <directory|files>",
	Short: "Generate data structures from .json/.js files",
	Run:   generateCmd.execute,
}

func init() {
	rootCmd.AddCommand(generateCmdCobra)
	cmd := generateCmdCobra
	cmd.Flags().BoolVarP(&generateCmd.verbose,
		"verbose", "v", false, "verbose output")
	cmd.Flags().StringVar(&generateCmd.pkg,
		"pkg", "", "package name of the generated code")
	cmd.Flags().StringVar(&generateCmd.outputDir,
		"outputDir", "", "the output directory")
	cmd.Flags().StringVar(&generateCmd.codeFileExt,
		"codeFileExt", "go", "filename extension of the generated code")
	cmd.Flags().StringVar(&generateCmd.tplDir,
		"tplDir", "", "the code template directory")
	cmd.Flags().BoolVar(&generateCmd.structCode,
		"x-struct", true, "generate struct related code")
	cmd.Flags().BoolVar(&generateCmd.structCollection,
		"x-collection", true, "generate collection related code")
	cmd.Flags().BoolVar(&generateCmd.collectionExtension,
		"x-collectionExtension", true, "generate collection extension related code")
	cmd.Flags().BoolVar(&generateCmd.easyjson,
		"x-easyjson", false, "generate easyjson related code")
	cmd.Flags().StringVar(&generateCmd.structNameSuffix,
		"structNameSuffix", "Conf", "name suffix of struct")

	generateCmd.registerSharedFlags(cmd)
}

type generateCmdT struct {
	verbose             bool
	pkg                 string
	outputDir           string
	codeFileExt         string
	tplDir              string
	structCode          bool
	structCollection    bool
	collectionExtension bool
	easyjson            bool
	structNameSuffix    string

	sharedFlags

	srcDir string
}

func (this *generateCmdT) execute(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		panic("no input file")
	}
	if this.outputDir == "" {
		panic("no --outputDir")
	}
	if this.codeFileExt == "" {
		panic("no --codeFileExt")
	} else if ext := strings.TrimSpace(this.codeFileExt); ext == "" {
		panic("invalid --codeFileExt")
	} else {
		this.codeFileExt = "." + strings.TrimLeft(ext, ".")
	}

	if this.tplDir != "" {
		loadTemplates(this.tplDir)
	}

	this.validateShared()

	allFiles := this.allFiles(args)
	if len(allFiles) == 0 {
		panic("no .json/.js file")
	}
	for _, file := range allFiles {
		var ext string
		switch {
		case strings.HasSuffix(file, ".json"):
			ext = ".json"
		case strings.HasSuffix(file, ".js"):
			ext = ".js"
		}
		if s := variable.ToCamel(strings.TrimSuffix(filepath.Base(file), ext)); s == "" {
			panic("invalid file name: " + file)
		}
	}
	for i := 0; i < len(allFiles); i++ {
		for j := i + 1; j < len(allFiles); j++ {
			si := filepath.Base(allFiles[i])
			sj := filepath.Base(allFiles[j])
			if si > sj {
				allFiles[i], allFiles[j] = allFiles[j], allFiles[i]
			}
		}
	}

	switch this.outputDir {
	case "":
		this.outputDir = filepath.Dir(allFiles[0])
	}
	if err := misc.FindDirectory(this.outputDir); err != nil {
		panic(fmt.Errorf("%s does not exist", this.outputDir))
	}
	switch this.pkg {
	case "":
		this.pkg = filepath.Base(filepath.Dir(allFiles[0]))
	}

	this.skipped = loadSkipped(allFiles)

	if this.verbose {
		fmt.Println("Output directory: " + this.outputDir)
		fmt.Println("Input files:\n\t" + strings.Join(allFiles, "\n\t"))
		fmt.Println()
	}

	if this.codeFileExt == ".go" {
		this.whichCmd("gofmt")
		this.whichCmd("easyjson")
	}

	p2 := filepath.Join(this.outputDir, TplCollectionExtension+this.codeFileExt)
	if err := misc.FindFile(p2); err == nil {
		if err := misc.FindFile(p2 + ".tmp"); err != nil {
			errRename := os.Rename(p2, p2+".tmp")
			if errRename != nil {
				panic(errRename)
			}
		}
	}
	restoreCollectionExtension := func() {
		if err := misc.FindFile(p2 + ".tmp"); err == nil {
			data, err := ioutil.ReadFile(p2 + ".tmp")
			if err != nil {
				fmt.Printf("Error: +%v\n", err)
			} else if err := ioutil.WriteFile(p2, data, 0644); err != nil {
				fmt.Printf("Error: +%v\n", err)
			} else if err := os.Remove(p2 + ".tmp"); err != nil {
				fmt.Printf("Error: +%v\n", err)
			}
		}
	}
	defer func() {
		restoreCollectionExtension()
	}()

	this.genStructRelatedCode(allFiles)
	this.genEasyJSONRelatedCode(allFiles)
}

func (this *generateCmdT) allFiles(args []string) []string {
	files1, err := misc.AllFiles(args, ".json", true)
	if err != nil {
		panic(err)
	}
	files2, err := misc.AllFiles(args, ".js", true)
	if err != nil {
		panic(err)
	}

	m2 := make(map[string]struct{})
	for _, f2 := range files2 {
		m2[f2] = struct{}{}
	}
	for _, f1 := range files1 {
		f2 := strings.TrimSuffix(f1, ".json") + ".js"
		if _, ok := m2[f2]; ok {
			panic(errors.Errorf("%s and %s cannot coexist under the same directory", f1, f2))
		}
	}

	return append(files1, files2...)
}

func deepToRef(node *guesser.Node) bool {
	for {
		switch node.ValueKind {
		case guesser.ValueKind_Primitive:
			return false
		case guesser.ValueKind_Struct:
			return false
		case guesser.ValueKind_Map:
			node = node.Value.MapValue
		case guesser.ValueKind_Array:
			node = node.Value.ArrayValue
		case guesser.ValueKind_Ref:
			return true
		case guesser.ValueKind_Undetermined:
			return false
		default:
			panic("impossible")
		}
	}
}

func deepToStruct(node *guesser.Node) bool {
	for {
		switch node.ValueKind {
		case guesser.ValueKind_Primitive:
			return false
		case guesser.ValueKind_Struct:
			return true
		case guesser.ValueKind_Map:
			node = node.Value.MapValue
		case guesser.ValueKind_Array:
			node = node.Value.ArrayValue
		case guesser.ValueKind_Ref:
			return false
		case guesser.ValueKind_Undetermined:
			return false
		default:
			panic("impossible")
		}
	}
}

func (this *generateCmdT) genStructRelatedCode(allFiles []string) {
	if !this.structCode {
		return
	}

	var jsonFiles []string
	var primaryStructNames []string
	var primaryStructNameMap = make(map[string]struct{})
	for _, file := range allFiles {
		basename1 := filepath.Base(file)
		var ext, jsonFile string
		switch {
		case strings.HasSuffix(basename1, ".json"):
			ext = ".json"
			jsonFile = basename1
		case strings.HasSuffix(basename1, ".js"):
			ext = ".js"
			jsonFile = strings.TrimSuffix(basename1, ".js") + ".json"
		default:
			panic("impossible")
		}
		basename2 := variable.ToCamel(strings.TrimSuffix(basename1, ext))
		jsonFiles = append(jsonFiles, jsonFile)
		primaryStructName := basename2 + this.structNameSuffix
		primaryStructNames = append(primaryStructNames, primaryStructName)
		primaryStructNameMap[primaryStructName] = struct{}{}
	}
	if len(jsonFiles) == 0 || len(primaryStructNames) == 0 {
		panic("impossible")
	}

	revRefGraph := make(map[string][]string)
	for i, file := range allFiles {
		fmt.Printf("Processing %s...\n", file)
		var g *guesser.Guesser
		switch {
		case strings.HasSuffix(file, ".json"):
			g = this.buildGuesser(file, primaryStructNameMap, this.structNameSuffix)
		case strings.HasSuffix(file, ".js"):
			if data, err := ioutil.ReadFile(file); err != nil {
				panic(err)
			} else {
				data = evalJavascript(data, file)
				g = this.buildGuesserWithJavascriptFile(
					data, file, primaryStructNameMap, this.structNameSuffix)
			}
		default:
			panic("impossible")
		}
		if g.Root.Meaningless() {
			panic("the file does not contain any meaningful data")
		}
		if this.verbose {
			fmt.Println(g.Root.Tree())
		}

		var newTypeNodes []*guesser.Node
		var newTypes = make(map[*guesser.Node]string)
		primaryStructName := primaryStructNames[i]
		newTypes[g.Root] = primaryStructName
		newTypeNodes = append(newTypeNodes, g.Root)

		structName := func(node *guesser.Node) string {
			if len(newTypes) == 1 {
				if g.Root.ValueKind == guesser.ValueKind_Map {
					if g.Root.Value.MapValue.ValueKind == guesser.ValueKind_Struct {
						for _, v := range newTypes {
							return v + "Item"
						}
					}
				}
			}
			n := chksum.Int31([]byte(node.Path()))
			return fmt.Sprintf("%s_%d", primaryStructName, n)
		}

		g.Root.RemoveUndeterminedChildren()
		g.Root.Traverse(func(node *guesser.Node) bool {
			switch node.ValueKind {
			case guesser.ValueKind_Primitive:
			case guesser.ValueKind_Struct:
				if _, ok := newTypes[node]; !ok {
					newTypes[node] = structName(node)
					newTypeNodes = append(newTypeNodes, node)
				}
			case guesser.ValueKind_Map:
			case guesser.ValueKind_Array:
			case guesser.ValueKind_Ref:
			case guesser.ValueKind_Undetermined:
			default:
				panic("impossible")
			}
			return true
		})

		tplArgs := struct {
			Pkg   string
			Nodes []*guesser.Node
		}{
			Pkg:   this.pkg,
			Nodes: newTypeNodes,
		}

		var funcMap = this.buildFuncMap(newTypes, jsonFiles[i])
		tpl, err := template.New(TplStruct).Funcs(funcMap).Parse(TplMap[TplStruct])
		if err != nil {
			panic(err)
		}

		var sb bytes.Buffer
		if err := tpl.Execute(&sb, tplArgs); err != nil {
			panic(err)
		}

		outputFile := filepath.Join(this.outputDir, primaryStructName+this.codeFileExt)
		err = ioutil.WriteFile(outputFile, sb.Bytes(), 0644)
		if err != nil {
			panic(err)
		}

		this.gofmt(outputFile)

		g.Root.Traverse(func(node *guesser.Node) bool {
			switch node.ValueKind {
			case guesser.ValueKind_Ref:
				key := node.Value.RawRef + ".json"
				idx := misc.IndexStrings(revRefGraph[key], jsonFiles[i])
				if idx < 0 {
					revRefGraph[key] = append(revRefGraph[key], jsonFiles[i])
				}
			}
			return true
		})
	}

	this.genCollectionFile(jsonFiles, primaryStructNames, revRefGraph)
}

func (this *generateCmdT) buildFuncMap(newTypes map[*guesser.Node]string, jsonFile string) map[string]interface{} {
	lookupStructName := func(node *guesser.Node) string {
		return newTypes[node]
	}

	var d int
	depth := func(action string) interface{} {
		switch action {
		case "+":
			d++
		case "-":
			d--
		case "0":
			d = 0
		case "v":
			return d
		default:
			panic(fmt.Errorf("unexpected action: %s", action))
		}
		return ""
	}

	var stack = deque.NewDeque()
	stackPush := func(str string) string {
		stack.PushBack(str)
		return ""
	}
	stackPop := func() string {
		return stack.PopBack().(string)
	}

	toPascal := func(name string) string {
		s := variable.ToPascal(name)
		if s == "" {
			return fmt.Sprintf("<!!!%s!!!>", name)
		}
		return s
	}
	toCamel := func(name string) string {
		s := variable.ToCamel(name)
		if s == "" {
			return fmt.Sprintf("<!!!%s!!!>", name)
		}
		return s
	}

	jsonFileFunc := func() string {
		return jsonFile
	}
	originalRefName := func(s string) string {
		return strings.TrimSuffix(s, this.structNameSuffix)
	}
	graveAccent := func() string {
		return "`"
	}

	return map[string]interface{}{
		"ucfirst":          misc.UCFirst,
		"lcfirst":          misc.LCFirst,
		"toPascal":         toPascal,
		"toCamel":          toCamel,
		"lookupStructName": lookupStructName,
		"deepToRef":        deepToRef,
		"deepToStruct":     deepToStruct,
		"stackPush":        stackPush,
		"stackPop":         stackPop,
		"depth":            depth,
		"trimPrefix":       strings.TrimPrefix,
		"trimSuffix":       strings.TrimSuffix,
		"originalRefName":  originalRefName,
		"hasPrefix":        strings.HasPrefix,
		"hasSuffix":        strings.HasSuffix,
		"jsonFile":         jsonFileFunc,
		"graveAccent":      graveAccent,
		"toUpper":          strings.ToUpper,
		"toLower":          strings.ToLower,
	}
}

func (this *generateCmdT) genCollectionFile(jsonFiles, primalStructNames []string, revRefGraph map[string][]string) {
	if !this.structCollection {
		return
	}

	for _, a := range revRefGraph {
		sort.Strings(a)
	}

	tplArgs := struct {
		Pkg                 string
		JSONFiles           []string
		Structs             []string
		RevRefGraph         map[string][]string
		CollectionExtension bool
	}{
		Pkg:                 this.pkg,
		JSONFiles:           jsonFiles,
		Structs:             primalStructNames,
		RevRefGraph:         revRefGraph,
		CollectionExtension: this.collectionExtension,
	}

	var funcMap = this.buildFuncMap(nil, "")
	tpl, err := template.New(TplCollection).Funcs(funcMap).Parse(TplMap[TplCollection])
	if err != nil {
		panic(err)
	}

	var sb bytes.Buffer
	if err := tpl.Execute(&sb, tplArgs); err != nil {
		panic(err)
	}

	outputFile := filepath.Join(this.outputDir, TplCollection+this.codeFileExt)
	err = ioutil.WriteFile(outputFile, sb.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

	this.gofmt(outputFile)
	this.genCollectionExtensionFile()
}

func (this *generateCmdT) genCollectionExtensionFile() {
	if !this.collectionExtension {
		return
	}

	tplArgs := struct {
		Pkg string
	}{
		Pkg: this.pkg,
	}

	var funcMap = this.buildFuncMap(nil, "")
	tpl, err := template.New(TplCollectionExtension).Funcs(funcMap).Parse(TplMap[TplCollectionExtension])
	if err != nil {
		panic(err)
	}

	var sb bytes.Buffer
	if err := tpl.Execute(&sb, &tplArgs); err != nil {
		panic(err)
	}

	outputFile := filepath.Join(this.outputDir, TplCollectionExtension+this.codeFileExt)
	if err := misc.FindFile(outputFile); err == nil {
		fmt.Printf("\n%s already exists, skipped.\n", outputFile)
		return
	}
	err = ioutil.WriteFile(outputFile, sb.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

	this.gofmt(outputFile)
}

func (this *generateCmdT) gofmt(file string) {
	if this.codeFileExt != ".go" {
		return
	}
	gofmt := exec.Command("gofmt", "-w", file)
	gofmt.Stdout = os.Stdout
	gofmt.Stderr = os.Stderr
	if err := gofmt.Run(); err != nil {
		panic(err)
	}
}

func (this *generateCmdT) genEasyJSONRelatedCode(allFiles []string) {
	if this.codeFileExt != ".go" {
		return
	}
	if !this.easyjson {
		return
	}

	var a1, a2 []string
	for _, file := range allFiles {
		fBase := filepath.Base(file)
		var extLen = 5
		if strings.HasSuffix(fBase, ".js") {
			extLen = 3
		}
		fName := variable.ToCamel(fBase[:len(fBase)-extLen]) + this.structNameSuffix + ".go"
		fPath := filepath.Join(this.outputDir, fName)
		a1 = append(a1, fPath)
		easyf := fPath[:len(fPath)-3] + easyjsonSuffix
		a2 = append(a2, easyf)
		_ = os.Remove(easyf)
	}

	fmt.Println()
	fmt.Println("Generating easyjson files...")
	this.genCodeWithEasyjson(a1...)
}

func (this *generateCmdT) whichCmd(name string) {
	if runtime.GOOS == "windows" {
		return
	}

	which := exec.Command("which", name)
	which.Stdout = ioutil.Discard
	which.Stderr = ioutil.Discard
	if err := which.Run(); err != nil {
		panic(fmt.Errorf("cannot find the following command: %s", name))
	}
}

func (this *generateCmdT) genCodeWithEasyjson(files ...string) {
	args := append([]string{"-all"}, files...)
	exePath := os.Args[0]
	absPath, err := filepath.Abs(exePath)
	if err != nil {
		panic(err)
	}
	easyjsonPath := filepath.Join(filepath.Dir(absPath), "easyjson")
	if err := misc.FindFile(easyjsonPath); err != nil {
		easyjsonPath = "easyjson"
	}
	easyjson := exec.Command(easyjsonPath, args...)
	easyjson.Stdout = os.Stdout
	easyjson.Stderr = os.Stderr
	if err := easyjson.Run(); err != nil {
		panic(err)
	}
}

func loadTemplates(dir string) {
	if err := misc.FindDirectory(dir); err != nil {
		panic(dir + " does not exist")
	}
	matches, err := filepath.Glob(filepath.Join(dir, "*.tpl"))
	if err != nil {
		panic(err)
	}

	var counter int
	defer func() {
		if counter > 0 {
			fmt.Println()
		}
	}()

	for _, m := range matches {
		abs, err := filepath.Abs(m)
		if err != nil {
			panic(err)
		}
		for k := range TplMap {
			if strings.HasSuffix(m, k+".tpl") {
				x, err := ioutil.ReadFile(abs)
				if err != nil {
					panic(err)
				}
				TplMap[k] = string(x)
				fmt.Printf("Loaded %s.tpl\n", k)
				counter++
			}
		}
	}
}
