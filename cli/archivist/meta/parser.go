package meta

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"github.com/kingsgroupos/misc"
	"github.com/kingsgroupos/misc/werr"
	"github.com/kingsgroupos/misc/wildcard"
)

type Parser struct {
	metaMap   map[string]*Meta
	wildMetas []*Meta
}

func NewParser() *Parser {
	return &Parser{
		metaMap: make(map[string]*Meta),
	}
}

func (this *Parser) ParseFile(file string) error {
	bts, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = this.ParseBytes(bts)
	if err != nil {
		return werr.NewRichError(err, "file", file)
	}
	return nil
}

func (this *Parser) ParseBytes(bts []byte) error {
	var lineNumber int
	scanner := bufio.NewScanner(bytes.NewReader(bts))
	for scanner.Scan() {
		lineNumber++
		meta, err := parseLine(scanner.Bytes())
		if err != nil {
			return werr.NewRichError(err, "line", lineNumber)
		}
		if meta == nil {
			continue
		}
		if _, ok := this.metaMap[meta.Path]; ok {
			return werr.NewRichError(errors.New("path duplication detected"), "line", lineNumber)
		}
		this.metaMap[meta.Path] = meta
		if strings.ContainsAny(meta.Path, "*?") {
			this.wildMetas = append(this.wildMetas, meta)
		}
	}

	return nil
}

func wildMatch(metas []*Meta, nodePath string) (*Meta, bool) {
	fields := misc.Split(nodePath, "/")
	numFields := len(fields)

outer:
	for _, mt := range metas {
		if len(mt.PathFields) != numFields {
			continue
		}
		for i := 0; i < numFields; i++ {
			if !wildcard.Match(mt.PathFields[i], fields[i]) {
				continue outer
			}
		}
		mt.Counter++
		return mt, true
	}

	return nil, false
}

func (this *Parser) Match(nodePath string) (*Meta, bool) {
	if mt, ok := this.metaMap[nodePath]; ok {
		mt.Counter++
		return mt, true
	}
	return nil, false
}

func (this *Parser) WildMatch(nodePath string) (*Meta, bool) {
	return wildMatch(this.wildMetas, nodePath)
}

func (this *Parser) UnusedMetas() []*Meta {
	var a []*Meta
	for _, mt := range this.metaMap {
		if mt.Counter == 0 {
			a = append(a, mt)
		}
	}
	return a
}
