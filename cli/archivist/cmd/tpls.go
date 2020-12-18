package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/kingsgroupos/archivist/cli/archivist/g"
	"github.com/kingsgroupos/misc"
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
