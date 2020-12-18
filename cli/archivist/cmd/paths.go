package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"gitlab-ee.funplus.io/watcher/watcher/archivist/cli/archivist/guesser"
)

var pathsCmd pathsCmdT

var pathsCmdCobra = &cobra.Command{
	Use:   "paths <file>",
	Short: "Show all node paths of a .json/.js file",
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
