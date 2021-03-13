package guesser

import (
	"io/ioutil"
	"strings"

	"github.com/kingsgroupos/misc"
)

var (
	suffixArr1 = [2]string{".json", ".js"}
	suffixArr2 = [2]string{".struct.json", ".struct.js"}
)

func ReadDataFile(filename string) ([]byte, error) {
	for i := 0; i < len(suffixArr1); i++ {
		suffix1 := suffixArr1[i]
		suffix2 := suffixArr2[i]
		if strings.HasSuffix(filename, suffix1) {
			if !strings.HasSuffix(filename, suffix2) {
				str1 := filename[:len(filename)-len(suffix1)]
				str2 := str1 + suffix2
				if misc.FindFile(str2) == nil {
					return ioutil.ReadFile(str2)
				}
			}
		}
	}

	return ioutil.ReadFile(filename)
}

func PickPureDataFiles(files []string) []string {
	a := make([]string, 0, len(files))
outer:
	for _, f := range files {
		for _, suffix := range suffixArr2 {
			if strings.HasSuffix(f, suffix) {
				continue outer
			}
		}
		a = append(a, f)
	}
	return a
}
