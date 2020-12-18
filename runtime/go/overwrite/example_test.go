package tmp

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/edwingeng/slog"
	"github.com/kingsgroupos/archivist/runtime/go/archivist"
	"github.com/kingsgroupos/archivist/runtime/go/overwrite/conf"
	"github.com/kingsgroupos/misc"
)

func Test_Overwrite(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(workingDir, "json")

	if err := os.Setenv(archivist.EnvName_ConfGroup, "develop"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv(archivist.EnvName_ConfSubgroup, "dolores"); err != nil {
		t.Fatal(err)
	}
	arch := archivist.NewArchivist(archivist.WithRoot(root), archivist.WithLogger(slog.NewDumbLogger()))
	c, err := arch.LoadCollection(conf.NewCollection)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	expected := conf.NewCollection().(*conf.Collection)
	expected.AConf = conf.AConf{
		A1: 100,
		A2: 200,
		A3: 300,
	}
	expected.BConf = conf.BConf{
		B1: 100,
		B2: 200,
		B3: 300,
	}
	expected.CConf = conf.CConf{
		C1: 1100,
		C2: 2200,
		C3: 1300,
	}
	expected.DConf = conf.DConf{
		D1: 100,
		D2: 1200,
		D3: 300,
	}
	expected.EConf = conf.EConf{
		E1: 100,
		E2: 200,
		E3: 1300,
		E5: []int64{1, 2, 3},
		E6: []int64{10, 20},
		E7: map[string]int64{
			"hello": 1,
			"world": 2,
		},
		E8: map[string]int64{
			"hello": 10,
		},
	}
	expected.FConf = conf.FConf{
		F1: 1100,
		F2: 1200,
		F3: 1300,
		F5: 1500,
	}
	expected.GConf = conf.GConf{
		"G1": 100,
		"G2": 1200,
		"G3": 300,
		"G5": 0,
	}
	expected.HConf = conf.HConf{
		"H1": []int64{10, 10, 10, 10, 10},
		"H2": []int64{20, 20},
		"H3": []int64{3, 3, 3},
		"H5": []int64{},
	}
	expected.IConf = conf.IConf{
		"I1": nil,
		"I2": map[string]int64{
			"i2.1": 2001,
		},
		"I3": map[string]int64{
			"i3.1": 31,
			"i3.2": 32,
		},
		"I5": map[string]int64{
			"i5.1": 51,
			"i5.2": 52,
			"i5.3": 53,
		},
	}
	expected.JConf = conf.JConf{
		"J1": "100",
		"J2": "1200",
		"J3": "300",
		"J5": "",
	}
	expected.KConf = conf.KConf{
		"K1": {
			Kxa: 11,
			Kxb: 128,
		},
		"K2": {
			Kxa: 21,
			Kxb: 220,
		},
		"K3": {
			Kxa: 31,
			Kxb: 328,
		},
		"K5": {
			Kxa: 51,
			Kxb: 52,
		},
	}

	for k, v := range c.Filename2Conf() {
		w := expected.Filename2Conf()[k]
		if !reflect.DeepEqual(v, w) {
			t.Fatalf("!reflect.DeepEqual(v, w). k: %s, v: %+v, w: %+v",
				k, misc.ToPrettyJSONString(v), misc.ToPrettyJSONString(w))
		}
	}

	err = archivist.Dump(c, filepath.Join(workingDir, "dump"))
	if err != nil {
		t.Fatal(err)
	}

	// in-memory
	var overwrites []archivist.Overwrite
	overwrites = append(overwrites,
		archivist.Overwrite{
			FileLevel: true,
			Target:    "a.json",
			Data:      []byte(`{}`),
		},
		archivist.Overwrite{
			FileLevel: true,
			Target:    "e.json",
			Data: []byte(`{
	"E7": {
		"a": 100,
		"b": 200
	}
}`),
		},
		archivist.Overwrite{
			FileLevel: false,
			Target:    "f.json",
			Data: []byte(`{
	"F1": 1100,
	"F2": 2000,
	"F3": 1300,
	"F5": 5000
}`),
		},
		archivist.Overwrite{
			FileLevel: false,
			Target:    "g.json",
			Data: []byte(`{
	"G5": 5000,
	"G9": 9000
}`),
		},
		archivist.Overwrite{
			FileLevel: false,
			Target:    "h.json",
			Data: []byte(`{
	"H5": null,
	"H9": [9, 9, 9]
}`),
		},
		archivist.Overwrite{
			FileLevel: false,
			Target:    "i.json",
			Data: []byte(`{
	"I5": {
		"i5.1": 511
	},
	"I9": {
		"i9.1": 91,
		"i9.2": 92
	}
}`),
		},
		archivist.Overwrite{
			FileLevel: false,
			Target:    "k.json",
			Data: []byte(`{
	"K5": {
		"Kxa": 511
	},
	"K9": {
		"Kxa": 91,
		"Kxb": 92
	}
}`),
		},
	)
	expected.AConf = conf.AConf{}
	expected.EConf = conf.EConf{
		E7: map[string]int64{
			"a": 100,
			"b": 200,
		},
	}
	expected.FConf = conf.FConf{
		F1: 1100,
		F2: 2000,
		F3: 1300,
		F5: 5000,
	}
	expected.GConf = conf.GConf{
		"G1": 100,
		"G2": 1200,
		"G3": 300,
		"G5": 5000,
		"G9": 9000,
	}
	expected.HConf = conf.HConf{
		"H1": []int64{10, 10, 10, 10, 10},
		"H2": []int64{20, 20},
		"H3": []int64{3, 3, 3},
		"H5": nil,
		"H9": []int64{9, 9, 9},
	}
	expected.IConf = conf.IConf{
		"I1": nil,
		"I2": map[string]int64{
			"i2.1": 2001,
		},
		"I3": map[string]int64{
			"i3.1": 31,
			"i3.2": 32,
		},
		"I5": map[string]int64{
			"i5.1": 511,
		},
		"I9": map[string]int64{
			"i9.1": 91,
			"i9.2": 92,
		},
	}
	expected.KConf = conf.KConf{
		"K1": {
			Kxa: 11,
			Kxb: 128,
		},
		"K2": {
			Kxa: 21,
			Kxb: 220,
		},
		"K3": {
			Kxa: 31,
			Kxb: 328,
		},
		"K5": {
			Kxa: 511,
			Kxb: 52,
		},
		"K9": {
			Kxa: 91,
			Kxb: 92,
		},
	}

	anotherC, err := arch.LoadCollection(conf.NewCollection, overwrites...)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	for k, v := range anotherC.Filename2Conf() {
		w := expected.Filename2Conf()[k]
		if !reflect.DeepEqual(v, w) {
			t.Fatalf("!reflect.DeepEqual(v, w). k: %s, v: %+v, w: %+v",
				k, misc.ToPrettyJSONString(v), misc.ToPrettyJSONString(w))
		}
	}

	patchedC, err := arch.PatchCollection(c, overwrites...)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	for k, v := range patchedC.Filename2Conf() {
		w := expected.Filename2Conf()[k]
		if !reflect.DeepEqual(v, w) {
			t.Fatalf("!reflect.DeepEqual(v, w). k: %s, v: %+v, w: %+v",
				k, misc.ToPrettyJSONString(v), misc.ToPrettyJSONString(w))
		}
	}
}
