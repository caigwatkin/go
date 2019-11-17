/*
Copyright 2018 Cai Gwatkin

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package secrets

import (
	"reflect"
	"testing"

	go_testing "github.com/caigwatkin/go/testing"
)

func TestReduceRequired(t *testing.T) {
	var data = []struct {
		desc     string
		input    []Required
		expected Required
	}{
		{
			desc:     "none",
			input:    nil,
			expected: make(Required),
		},
		{
			desc:     "empty",
			input:    []Required{},
			expected: make(Required),
		},
		{
			desc: "single set",
			input: []Required{
				Required{
					"111": []string{"111"},
				},
			},
			expected: Required{
				"111": []string{"111"},
			},
		},
		{
			desc: "single set with duplicate value",
			input: []Required{
				Required{
					"111": []string{"111", "111"},
				},
			},
			expected: Required{
				"111": []string{"111"},
			},
		},
		{
			desc: "two small sets",
			input: []Required{
				Required{
					"111": []string{"111"},
				},
				Required{
					"222": []string{"222"},
				},
			},
			expected: Required{
				"111": []string{"111"},
				"222": []string{"222"},
			},
		},
		{
			desc: "three small sets",
			input: []Required{
				Required{
					"111": []string{"111"},
				},
				Required{
					"222": []string{"222"},
				},
				Required{
					"333": []string{"333"},
				},
			},
			expected: Required{
				"111": []string{"111"},
				"222": []string{"222"},
				"333": []string{"333"},
			},
		},
		{
			desc: "two small sets, multiple values in string arrays",
			input: []Required{
				Required{
					"111": []string{"111", "1"},
				},
				Required{
					"222": []string{"222", "2"},
				},
			},
			expected: Required{
				"111": []string{"111", "1"},
				"222": []string{"222", "2"},
			},
		},
		{
			desc: "two small sets, same key",
			input: []Required{
				Required{
					"111": []string{"111"},
				},
				Required{
					"111": []string{"1111"},
				},
			},
			expected: Required{
				"111": []string{"111", "1111"},
			},
		},
		{
			desc: "two small sets, same key, same value",
			input: []Required{
				Required{
					"111": []string{"111"},
				},
				Required{
					"111": []string{"111"},
				},
			},
			expected: Required{
				"111": []string{"111"},
			},
		},
		{
			desc: "two small sets, empty string key",
			input: []Required{
				Required{
					"": []string{""},
				},
				Required{
					"111": []string{"1111"},
				},
			},
			expected: Required{
				"":    []string{""},
				"111": []string{"1111"},
			},
		},
		{
			desc: "two small sets, first one empty",
			input: []Required{
				Required{},
				Required{
					"111": []string{"111"},
				},
			},
			expected: Required{
				"111": []string{"111"},
			},
		},
		{
			desc: "two small sets, second one empty",
			input: []Required{
				Required{
					"111": []string{"111"},
				},
				Required{},
			},
			expected: Required{
				"111": []string{"111"},
			},
		},
	}

	for i, d := range data {
		result := ReduceRequired(d.input...)
		if !reflect.DeepEqual(result, d.expected) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected,
				Result:     result,
			}))
		}
	}
}
