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

package testing

import (
	"testing"
)

func TestErrorf(t *testing.T) {

	var data = []struct {
		desc     string
		input    Error
		expected string
	}{
		{
			desc: "all properties with values",
			input: Error{
				Unexpected: "a",
				Desc:       "b",
				At:         1,
				Input:      "c",
				Expected:   "d",
				Result:     "e",
			},
			expected: `{
	"Unexpected": "a",
	"Desc": "b",
	"At": 1,
	"Input": "c",
	"Expected": "d",
	"Result": "e"
}`,
		},

		{
			desc: "all properties with zero-values",
			input: Error{
				Unexpected: "",
				Desc:       "",
				At:         0,
				Input:      nil,
				Expected:   nil,
				Result:     nil,
			},
			expected: `{
	"Unexpected": "",
	"Desc": "",
	"At": 0,
	"Input": null,
	"Expected": null,
	"Result": null
}`,
		},

		{
			desc: "not json marshallable interfaces",
			input: Error{
				Unexpected: "",
				Desc:       "",
				At:         0,
				Input:      func() {},
				Expected:   func() {},
				Result:     func() {},
			},
			expected: `{
	"Unexpected": "",
	"Desc": "",
	"At": 0,
	"Input": "potentially unmarshallable",
	"Expected": "potentially unmarshallable",
	"Result": "potentially unmarshallable"
}`,
		},
	}

	for i, d := range data {
		result := Errorf(d.input)

		if result != d.expected {
			t.Errorf(`{
	"Unexpected": %q,
	"Desc": %q,
	"At": %d,
	"Input": %v,
	"Expected": %v,
	"Result": %v
}`, "result", d.desc, i, d.input, d.expected, result)
		}
	}
}
