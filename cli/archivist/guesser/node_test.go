// BSD 3-Clause License
//
// Copyright (c) 2020, Kingsgroup
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package guesser

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func TestNode_parseValue(t *testing.T) {
	cases := []struct {
		jsonText      string
		expectedErr   error
		expectedPaths []string
	}{
		{
			jsonText:    `{}`,
			expectedErr: ErrSkip,
		},
		{
			jsonText: `[]`,
			expectedPaths: []string{
				"/", "[]",
				"/[]", "undetermined",
			},
		},
		{
			jsonText:    `[{}]`,
			expectedErr: ErrSkip,
		},
		{
			jsonText:    `[{},{}]`,
			expectedErr: ErrSkip,
		},
		{
			jsonText: `[[]]`,
			expectedPaths: []string{
				"/", "[]",
				"/[]", "[]",
				"/[]/[]", "undetermined",
			},
		},
		{
			jsonText: `[[],[]]`,
			expectedPaths: []string{
				"/", "[]",
				"/[]", "[]",
				"/[]/[]", "undetermined",
			},
		},
		{
			jsonText: `[[[]]]`,
			expectedPaths: []string{
				"/", "[]",
				"/[]", "[]",
				"/[]/[]", "[]",
				"/[]/[]/[]", "undetermined",
			},
		},
		{
			jsonText:    `[1]`,
			expectedErr: nil,
			expectedPaths: []string{
				"/", "[]",
				"/[]", "int64",
			},
		},
		{
			jsonText: `[1, 2, 3]`,
			expectedPaths: []string{
				"/", "[]",
				"/[]", "int64",
			},
		},
		{
			jsonText: `[1.5]`,
			expectedPaths: []string{
				"/", "[]",
				"/[]", "float64",
			},
		},
		{
			jsonText: `[1, 2, 3.5]`,
			expectedPaths: []string{
				"/", "[]",
				"/[]", "float64",
			},
		},
		{
			jsonText: `["hello"]`,
			expectedPaths: []string{
				"/", "[]",
				"/[]", "string",
			},
		},
		{
			jsonText: `["hello", "world"]`,
			expectedPaths: []string{
				"/", "[]",
				"/[]", "string",
			},
		},
		{
			jsonText: `[true]`,
			expectedPaths: []string{
				"/", "[]",
				"/[]", "bool",
			},
		},
		{
			jsonText: `[true, false]`,
			expectedPaths: []string{
				"/", "[]",
				"/[]", "bool",
			},
		},
		{
			jsonText: `[[1], [2, 3]]`,
			expectedPaths: []string{
				"/", "[]",
				"/[]", "[]",
				"/[]/[]", "int64",
			},
		},
		{
			jsonText: `[[1], [2, 3.5]]`,
			expectedPaths: []string{
				"/", "[]",
				"/[]", "[]",
				"/[]/[]", "float64",
			},
		},
		{
			jsonText: `[["hello"], ["world"]]`,
			expectedPaths: []string{
				"/", "[]",
				"/[]", "[]",
				"/[]/[]", "string",
			},
		},
		{
			jsonText: `[[true], [false]]`,
			expectedPaths: []string{
				"/", "[]",
				"/[]", "[]",
				"/[]/[]", "bool",
			},
		},
		{
			jsonText: `
{
	"A": {
		"B": {
			"C": 100
		}
	}
}
`,
			expectedPaths: []string{
				"/", "{}",
				"/A", "{}",
				"/A/B", "{}",
				"/A/B/C", "int64",
			},
		},
		{
			jsonText: `
[
	{
		"A": "hello",
		"B": true
	},
	{
		"C": 100,
		"D": 1.5,
		"E": ["x", "y", "z"],
		"F": [
			{"X": 100},
			{"B": false}
		]
	}
]
`,
			expectedPaths: []string{
				"/", "[]",
				"/[]", "{}",
				"/[]/A", "string",
				"/[]/B", "bool",
				"/[]/C", "int64",
				"/[]/D", "float64",
				"/[]/E", "[]",
				"/[]/E/[]", "string",
				"/[]/F", "[]",
				"/[]/F/[]", "{}",
				"/[]/F/[]/B", "bool",
				"/[]/F/[]/X", "int64",
			},
		},
		{
			jsonText: `
[
	{
		"A": 100
	},
	{
		"A": "hello"
	}
]
`,
			expectedErr: ErrIncompatible,
		},
		{
			jsonText: `
[
	{
		"A": 100
	},
	{
		"A": [1, 2, 3]
	}
]
`,
			expectedErr: ErrIncompatible,
		},
		{
			jsonText: `
{
	"30": {
		"A": 100
	},
	"50": {
		"A": false
	}
}
`,
			expectedErr: ErrIncompatible,
		},
	}

	for i, c := range cases {
		var obj interface{}
		if err := json.Unmarshal([]byte(c.jsonText), &obj); err != nil {
			t.Fatalf("failed to parse the json text. i: %d, err: %s\n%s", i, err, c.jsonText)
		}
		root := NewNode("", nil)
		if err := root.parseValue(obj, nil); err != nil {
			switch c.expectedErr {
			case ErrSkip:
				if errors.Is(err, ErrSkip) {
					continue
				}
				t.Fatalf("parseValue failed [1]. i: %d, err: %s\n%s", i, err, c.jsonText)
			case ErrIncompatible:
				if errors.Is(err, ErrIncompatible) {
					continue
				}
				fmt.Print(PrettyErrorIncompatible(err))
				t.Fatalf("parseValue failed [2]. i: %d, err: %s\n%s", i, err, c.jsonText)
			default:
				t.Fatalf("parseValue failed [3]. i: %d, err: %s\n%s", i, err, c.jsonText)
			}
		} else {
			if c.expectedErr != nil {
				t.Fatalf("parseValue should return an error. i: %d, excepted: %s\n%s",
					i, c.expectedErr, c.jsonText)
			}
		}

		a := root.AllPathTypes()
		if len(a) != len(c.expectedPaths) {
			t.Fatalf("len(a) != len(c.expectedPaths). i: %d, n1: %d, n2: %d", i, len(a), len(c.expectedPaths))
		}
		for j := 0; j < len(a); j += 2 {
			if a[j] != c.expectedPaths[j] {
				t.Fatalf("a[j] != c.expectedPaths[j]. i: %d, j: %d, a[j]: %q, c.expectedPaths[j]: %q",
					i, j, a[j], c.expectedPaths[j])
			}
			if a[j+1] != c.expectedPaths[j+1] {
				t.Fatalf("a[j+1] != c.expectedPaths[j+1]. i: %d, j: %d, path: %q, a[j+1]: %q, c.expectedPaths[j+1]: %q",
					i, j, a[j], a[j+1], c.expectedPaths[j+1])
			}
		}
	}
}

func TestNode_Path(t *testing.T) {
	g := NewGuesser()
	err := g.ParseFile(`../example/json/workbook1.json`)
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{
		"/", "map[int64]",
		"/map[]", "{}",
		"/map[]/ID", "int64",
		"/map[]/MyBool", "bool",
		"/map[]/MyBoolArray", "[]",
		"/map[]/MyBoolArray/[]", "bool",
		"/map[]/MyDatetime", "string",
		"/map[]/MyDatetimeArray", "[]",
		"/map[]/MyDatetimeArray/[]", "string",
		"/map[]/MyDuration", "string",
		"/map[]/MyDurationArray", "[]",
		"/map[]/MyDurationArray/[]", "string",
		"/map[]/MyFloat", "float64",
		"/map[]/MyFloatArray", "[]",
		"/map[]/MyFloatArray/[]", "float64",
		"/map[]/MyI18N", "string",
		"/map[]/MyI18NArray", "[]",
		"/map[]/MyI18NArray/[]", "string",
		"/map[]/MyInt", "int64",
		"/map[]/MyIntArray", "[]",
		"/map[]/MyIntArray/[]", "int64",
		"/map[]/MyRef", "int64",
		"/map[]/MyRefArray", "[]",
		"/map[]/MyRefArray/[]", "int64",
		"/map[]/MyString", "string",
		"/map[]/MyStringArray", "[]",
		"/map[]/MyStringArray/[]", "string",
	}

	a := g.Root.AllPathTypes()
	if len(a) != len(expected) {
		t.Fatalf("len(a) != len(expected). n1: %d, n2: %d", len(a), len(expected))
	}
	for i := 0; i < len(a); i += 2 {
		if a[i] != expected[i] {
			t.Fatalf("a[i] != expected[i]. i: %d, a[i]: %q, expected[i]: %q", i, a[i], expected[i])
		}
		if a[i+1] != expected[i+1] {
			t.Fatalf("a[i+1] != expected[i+1]. i: %d, path: %q, a[i+1]: %q, expected[i+1]: %q",
				i, a[i], a[i+1], expected[i+1])
		}
	}
}
