package archivist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gitlab-ee.funplus.io/watcher/watcher/misc"
)

const EnvName_ConfGroup = "WATCHER_CONF_GROUP"
const EnvName_ConfSubgroup = "WATCHER_CONF_SUBGROUP"

var (
	errEnv = fmt.Errorf("environmental variable %s should not be empty", EnvName_ConfGroup)
)

func ConfGroup() string {
	return os.Getenv(EnvName_ConfGroup)
}

func ConfSubgroup() string {
	return os.Getenv(EnvName_ConfSubgroup)
}

func MustHaveConfGroup() {
	if ConfGroup() == "" {
		panic(errEnv)
	}
}

func Dump(c Collection, outputDir string) error {
	if c == nil {
		return errors.New("c cannot be nil")
	}
	if err := misc.FindDirectory(outputDir); err != nil {
		return errors.WithStack(err)
	}

	for k, v := range c.Filename2Conf() {
		if v == nil {
			continue
		}
		data1, err := json.Marshal(v)
		if err != nil {
			return errors.Errorf("json.Marshal failed. file: %s, err: %+v", k, err)
		}
		var obj interface{}
		if err := json.Unmarshal(data1, &obj); err != nil {
			return errors.Errorf("json.Unmarshal failed. file: %s, err: %+v", k, err)
		}
		data2, err := json.MarshalIndent(obj, "", "\t")
		if err != nil {
			return errors.Errorf("json.MarshalIndent failed. file: %s, err: %+v", k, err)
		}

		data2 = append(data2, '\n')
		file := filepath.Join(outputDir, k)
		if err := ioutil.WriteFile(file, data2, 0644); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
