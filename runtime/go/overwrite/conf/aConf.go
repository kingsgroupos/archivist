// Code generated by archivist. DO NOT EDIT.

package conf

import (
	"time"

	"github.com/pkg/errors"
	"gitlab-ee.funplus.io/watcher/watcher/archivist/runtime/go/archivist"
	"gitlab-ee.funplus.io/watcher/watcher/misc/wtime"
)

var (
	_ = time.After
	_ = errors.New
	_ = archivist.NewArchivist
	_ = wtime.ParseDuration
)

// easyjson:json
type AConf struct {
	A1 int64 `json:"A1" bson:"A1"`
	A2 int64 `json:"A2" bson:"A2"`
	A3 int64 `json:"A3" bson:"A3"`
}

func (this *AConf) bindRefs(c *Collection) error {
	if this == nil {
		return nil
	}

	var ok bool
	_ = ok

	return nil
}