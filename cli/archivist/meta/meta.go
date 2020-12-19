package meta

import (
	"bytes"
	"strings"

	"github.com/kingsgroupos/misc"
	"github.com/pkg/errors"
)

type Meta struct {
	Path    string
	NewType string
	Notes   string
	RawData []byte
	Counter int

	PathFields []string
}

var (
	validTypes = map[string]struct{}{
		"bool":     {},
		"int":      {},
		"int8":     {},
		"int16":    {},
		"int32":    {},
		"int64":    {},
		"uint":     {},
		"uint8":    {},
		"uint16":   {},
		"uint32":   {},
		"uint64":   {},
		"float32":  {},
		"float64":  {},
		"string":   {},
		"datetime": {},
		"duration": {},
		"i18n":     {},
		"[]":       {},
		"{}":       {},
	}
)

func parseLine(line []byte) (*Meta, error) {
	trimmed := bytes.TrimSpace(line)
	if len(trimmed) == 0 {
		return nil, nil
	}
	fields := bytes.Fields(line)
	if len(fields) < 2 {
		return nil, errors.New("too few fields")
	}
	if len(fields) > 3 {
		fields[2] = bytes.Join(fields[2:], []byte(" "))
		fields = fields[:3]
	}

	mt := &Meta{
		Path:    string(fields[0]),
		NewType: string(fields[1]),
		RawData: append([]byte{}, line...),
	}
	if len(fields) > 2 {
		mt.Notes = string(fields[2])
	}
	if !isValidNewType(mt.NewType) {
		return nil, errors.Errorf("invalid new type: %s. path: %s", mt.NewType, mt.Path)
	}

	mt.PathFields = misc.Split(mt.Path, "/")
	return mt, nil
}

func isValidNewType(t string) bool {
	if _, ok := validTypes[t]; ok {
		return true
	} else if strings.HasPrefix(t, "map[") && strings.HasSuffix(t, "]") {
		str := t[4 : len(t)-1]
		if _, ok := validTypes[str]; ok {
			return true
		}
	} else if strings.HasPrefix(t, "ref@") {
		return true
	}
	return false
}
