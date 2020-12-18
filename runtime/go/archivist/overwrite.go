package archivist

import (
	"sort"

	"github.com/pkg/errors"
)

type Overwrite struct {
	FileLevel bool
	Target    string
	Data      []byte
}

type organizedOverwrites struct {
	fileLevel    map[string][][]byte
	contentLevel map[string][][]byte
}

func (this *organizedOverwrites) organize(c Collection, o []Overwrite) error {
	if len(o) == 0 {
		return nil
	}

	m1 := make(map[string][][]byte)
	m2 := make(map[string][][]byte)
	targets := c.Filename2Conf()
	for _, x := range o {
		if _, ok := targets[x.Target]; !ok {
			return errors.Errorf("<archivist> unknown overwrite target: %s", x.Target)
		}
		if x.FileLevel {
			m1[x.Target] = append(m1[x.Target], x.Data)
		} else {
			m2[x.Target] = append(m2[x.Target], x.Data)
		}
	}

	this.fileLevel = m1
	this.contentLevel = m2
	return nil
}

func (this *organizedOverwrites) findAffected(revRefGraph map[string][]string) []string {
	var pending []string
	for f := range this.fileLevel {
		pending = append(pending, f)
	}
	for f := range this.contentLevel {
		pending = append(pending, f)
	}

	m := make(map[string]struct{})
	affected := make(map[string]struct{})
	for n := len(pending); n > 0; n = len(pending) {
		f := pending[n-1]
		pending = pending[:n-1]
		if _, ok := m[f]; ok {
			continue
		}

		m[f] = struct{}{}
		for _, referer := range revRefGraph[f] {
			affected[referer] = struct{}{}
			pending = append(pending, referer)
		}
	}

	for f := range this.fileLevel {
		delete(affected, f)
	}
	for f := range this.contentLevel {
		delete(affected, f)
	}

	a := make([]string, 0, len(affected))
	for f := range affected {
		a = append(a, f)
	}
	sort.Strings(a)
	return a
}
