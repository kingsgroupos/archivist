package guesser

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type Guesser struct {
	Root            *Node
	stringKeyTester func(nodePath string) bool
}

func NewGuesser(opts ...Option) *Guesser {
	g := &Guesser{}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

func (this *Guesser) ParseFile(file string) error {
	bts, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return this.ParseBytes(bts)
}

func (this *Guesser) ParseBytes(bts []byte) error {
	var v interface{}
	if err := json.Unmarshal(bts, &v); err != nil {
		return err
	}
	this.Root = NewNode("", nil)
	err := this.Root.parseValue(v, this.stringKeyTester)
	if errors.Is(err, ErrSkip) {
		return errors.Errorf("it is impossible deduce the data structure of a file having no data")
	}
	return err
}

type Option func(g *Guesser)

func WithStringKeyTester(f func(nodePath string) bool) Option {
	return func(g *Guesser) {
		g.stringKeyTester = f
	}
}
