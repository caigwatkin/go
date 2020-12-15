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

package log

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"

	go_environment "github.com/caigwatkin/go/environment"
	go_errors "github.com/caigwatkin/go/errors"
	go_testing "github.com/caigwatkin/go/testing"
)

func Test_NewClient(t *testing.T) {
	var data = []struct {
		desc     string
		input    Config
		expected client
	}{
		{
			desc: "config",
			input: Config{
				Env: go_environment.Environment{
					App: "app",
				},
			},
			expected: client{
				config: Config{
					Env: go_environment.Environment{
						App: "app",
					},
				},
			},
		},
	}

	for i, d := range data {
		result := NewClient(context.Background(), d.input)

		if reflect.TypeOf(result) != reflect.TypeOf(d.expected) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   reflect.TypeOf(d.expected),
				Result:     reflect.TypeOf(result),
			}))
		}
		if v, ok := result.(client); !ok {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected,
				Result:     result,
			}))

		} else if v.config != d.expected.config {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.config",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.config,
				Result:     v.config,
			}))
		}
	}
}

func Test_FmtAny(t *testing.T) {
	notJSONMarshallableFunc := func() {}
	type valueStruct struct {
		X string `json:"x,omitempty"`
		Y int    `json:"y,omitempty"`
		Z bool   `json:"z,omitempty"`
	}
	type input struct {
		Value interface{}
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc  string
		input input
		expected
	}{
		{
			desc: "struct",
			input: input{
				Value: valueStruct{
					X: "some_string",
					Y: 2,
					Z: true,
				},
				Name: "name",
			},
			expected: expected{
				Result:       "\"name\": {\n\t\t\"type\": \"log.valueStruct\",\n\t\t\"value\": {\n\t\t\t\"x\": \"some_string\",\n\t\t\t\"y\": 2,\n\t\t\t\"z\": true\n\t\t}\n\t}",
				ResultRemote: "\"name\":{\"type\":\"log.valueStruct\",\"value\":{\"x\":\"some_string\",\"y\":2,\"z\":true}}",
			},
		},

		{
			desc: "struct omitempty",
			input: input{
				Value: valueStruct{
					X: "",
					Y: 0,
					Z: false,
				},
				Name: "name",
			},
			expected: expected{
				Result:       "\"name\": {\n\t\t\"type\": \"log.valueStruct\",\n\t\t\"value\": {}\n\t}",
				ResultRemote: "\"name\":{\"type\":\"log.valueStruct\",\"value\":{}}",
			},
		},

		{
			desc: "not JSON marshallable",
			input: input{
				Value: notJSONMarshallableFunc,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": {\n\t\t\"type\": \"func()\",\n\t\t\"value\": \"NOT JSON MARSHALLABLE\"\n\t}",
				ResultRemote: "\"name\":{\"type\":\"func()\",\"value\":\"NOT JSON MARSHALLABLE\"}",
			},
		},

		{
			desc: "nil",
			input: input{
				Value: nil,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": null",
				ResultRemote: "\"name\":null",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtAny(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtAny(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtAnys(t *testing.T) {
	notJSONMarshallableFunc := func() {}
	type input struct {
		Value []interface{}
		Name  string
	}
	type valueStruct struct {
		X string `json:"x,omitempty"`
		Y int    `json:"y,omitempty"`
		Z bool   `json:"z,omitempty"`
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "single",
			input: input{
				Value: []interface{}{
					valueStruct{
						X: "some_string",
						Y: 2,
						Z: true,
					},
				},
				Name: "name",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\t{\n\t\t\t\"type\": \"log.valueStruct\",\n\t\t\t\"value\": {\n\t\t\t\t\"x\": \"some_string\",\n\t\t\t\t\"y\": 2,\n\t\t\t\t\"z\": true\n\t\t\t}\n\t\t}\n\t]",
				ResultRemote: "\"name\":[{\"type\":\"log.valueStruct\",\"value\":{\"x\":\"some_string\",\"y\":2,\"z\":true}}]",
			},
		},

		{
			desc: "multi",
			input: input{
				Value: []interface{}{
					valueStruct{
						X: "some_string",
						Y: 2,
						Z: true,
					},
					nil,
					notJSONMarshallableFunc,
				},
				Name: "name",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\t{\n\t\t\t\"type\": \"log.valueStruct\",\n\t\t\t\"value\": {\n\t\t\t\t\"x\": \"some_string\",\n\t\t\t\t\"y\": 2,\n\t\t\t\t\"z\": true\n\t\t\t}\n\t\t},\n\t\tnull,\n\t\t{\n\t\t\t\"type\": \"func()\",\n\t\t\t\"value\": \"NOT JSON MARSHALLABLE\"\n\t\t}\n\t]",
				ResultRemote: "\"name\":[{\"type\":\"log.valueStruct\",\"value\":{\"x\":\"some_string\",\"y\":2,\"z\":true}},null,{\"type\":\"func()\",\"value\":\"NOT JSON MARSHALLABLE\"}]",
			},
		},

		{
			desc: "empty",
			input: input{
				Value: []interface{}{},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},

		{
			desc: "nil",
			input: input{
				Value: nil,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtAnys(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtAnys(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtBool(t *testing.T) {
	type input struct {
		Value bool
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "false",
			input: input{
				Value: false,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": false",
				ResultRemote: "\"name\":false",
			},
		},

		{
			desc: "true",
			input: input{
				Value: true,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": true",
				ResultRemote: "\"name\":true",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtBool(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtBool(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtBools(t *testing.T) {
	type input struct {
		Value []bool
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "single",
			input: input{
				Value: []bool{true},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\ttrue\n\t]",
				ResultRemote: "\"name\":[true]",
			},
		},

		{
			desc: "multi",
			input: input{
				Value: []bool{true, false},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\ttrue,\n\t\tfalse\n\t]",
				ResultRemote: "\"name\":[true,false]",
			},
		},

		{
			desc: "empty",
			input: input{
				Value: []bool{},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},

		{
			desc: "nil",
			input: input{
				Value: nil,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtBools(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtBools(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtByte(t *testing.T) {
	type input struct {
		Value byte
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "Character",
			input: input{
				Value: 'A',
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": 'A'",
				ResultRemote: "\"name\":'A'",
			},
		},

		{
			desc: "integer",
			input: input{
				Value: 0,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": '\\x00'",
				ResultRemote: "\"name\":'\\x00'",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtByte(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtByte(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtBytes(t *testing.T) {
	type input struct {
		Value []byte
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "nil",
			input: input{
				Value: nil,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": \"\"",
				ResultRemote: "\"name\":\"\"",
			},
		},

		{
			desc: "empty",
			input: input{
				Value: []byte{},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": \"\"",
				ResultRemote: "\"name\":\"\"",
			},
		},

		{
			desc: "string",
			input: input{
				Value: []byte("some_string"),
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": \"some_string\"",
				ResultRemote: "\"name\":\"some_string\"",
			},
		},

		{
			desc: "characters",
			input: input{
				Value: []byte{'A', 'a', 'B'},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": \"AaB\"",
				ResultRemote: "\"name\":\"AaB\"",
			},
		},

		{
			desc: "characters",
			input: input{
				Value: []byte{0, 25, 18},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": \"\\x00\\x19\\x12\"",
				ResultRemote: "\"name\":\"\\x00\\x19\\x12\"",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtBytes(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtBytes(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtDuration(t *testing.T) {
	type input struct {
		Value time.Duration
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "nano",
			input: input{
				Value: time.Nanosecond,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": \"1ns\"",
				ResultRemote: "\"name\":\"1ns\"",
			},
		},

		{
			desc: "micro",
			input: input{
				Value: time.Microsecond,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": \"1µs\"",
				ResultRemote: "\"name\":\"1µs\"",
			},
		},

		{
			desc: "milli",
			input: input{
				Value: time.Millisecond,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": \"1ms\"",
				ResultRemote: "\"name\":\"1ms\"",
			},
		},

		{
			desc: "second",
			input: input{
				Value: time.Second,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": \"1s\"",
				ResultRemote: "\"name\":\"1s\"",
			},
		},

		{
			desc: "minute",
			input: input{
				Value: time.Minute,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": \"1m0s\"",
				ResultRemote: "\"name\":\"1m0s\"",
			},
		},

		{
			desc: "hour",
			input: input{
				Value: time.Hour,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": \"1h0m0s\"",
				ResultRemote: "\"name\":\"1h0m0s\"",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtDuration(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtDuration(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtDurations(t *testing.T) {
	type input struct {
		Value []time.Duration
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "single",
			input: input{
				Value: []time.Duration{time.Nanosecond},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\t\"1ns\"\n\t]",
				ResultRemote: "\"name\":[\"1ns\"]",
			},
		},

		{
			desc: "multi",
			input: input{
				Value: []time.Duration{time.Nanosecond, time.Microsecond},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\t\"1ns\",\n\t\t\"1µs\"\n\t]",
				ResultRemote: "\"name\":[\"1ns\",\"1µs\"]",
			},
		},

		{
			desc: "empty",
			input: input{
				Value: []time.Duration{},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},

		{
			desc: "nil",
			input: input{
				Value: nil,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtDurations(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtDurations(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtError(t *testing.T) {
	errWithTrace := go_errors.New("error")
	trace := fmt.Sprintf("%+v", errWithTrace)
	type input struct {
		err error
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "nil",
			input: input{
				err: nil,
			},
			expected: expected{
				Result:       "\"error\": null",
				ResultRemote: "\"error\":null",
			},
		},

		{
			desc: "no trace",
			input: input{
				err: errors.New("some_string"),
			},
			expected: expected{
				Result:       "\"error\": \"some_string\"",
				ResultRemote: "\"error\":\"some_string\"",
			},
		},

		{
			desc: "trace",
			input: input{
				err: errWithTrace,
			},
			expected: expected{
				Result:       fmt.Sprintf("\"error\": {\n\t\t\"friendly\": \"error\",\n\t\t\"trace\": %s\n\t}", trace),
				ResultRemote: fmt.Sprintf("\"error\":{\"friendly\":\"error\",\"trace\":%s}", trace),
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtError(d.input.err)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtError(d.input.err)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtFloat32(t *testing.T) {
	type input struct {
		Value float32
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "zero",
			input: input{
				Value: 0,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": 0.00000",
				ResultRemote: "\"name\":0.00000",
			},
		},

		{
			desc: "positive whole",
			input: input{
				Value: 1,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": 1.00000",
				ResultRemote: "\"name\":1.00000",
			},
		},

		{
			desc: "positive with dp",
			input: input{
				Value: 1.12345,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": 1.12345",
				ResultRemote: "\"name\":1.12345",
			},
		},

		{
			desc: "positive with dp greater than 5 round down",
			input: input{
				Value: 1.123454,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": 1.12345",
				ResultRemote: "\"name\":1.12345",
			},
		},

		{
			desc: "positive with dp greater than 5 round up",
			input: input{
				Value: 1.123455,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": 1.12346",
				ResultRemote: "\"name\":1.12346",
			},
		},

		{
			desc: "negative whole",
			input: input{
				Value: -1,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": -1.00000",
				ResultRemote: "\"name\":-1.00000",
			},
		},

		{
			desc: "negative with dp",
			input: input{
				Value: -1.12345,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": -1.12345",
				ResultRemote: "\"name\":-1.12345",
			},
		},

		{
			desc: "negative with dp greater than 5 round down",
			input: input{
				Value: -1.123454,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": -1.12345",
				ResultRemote: "\"name\":-1.12345",
			},
		},

		{
			desc: "negative with dp greater than 5 round up",
			input: input{
				Value: -1.123455,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": -1.12346",
				ResultRemote: "\"name\":-1.12346",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtFloat32(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtFloat32(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtFloat32s(t *testing.T) {
	type input struct {
		Value []float32
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "single",
			input: input{
				Value: []float32{1.2},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\t1.20000\n\t]",
				ResultRemote: "\"name\":[1.20000]",
			},
		},

		{
			desc: "multi",
			input: input{
				Value: []float32{1.2, 3.4},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\t1.20000,\n\t\t3.40000\n\t]",
				ResultRemote: "\"name\":[1.20000,3.40000]",
			},
		},

		{
			desc: "empty",
			input: input{
				Value: []float32{},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},

		{
			desc: "nil",
			input: input{
				Value: nil,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtFloat32s(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtFloat32s(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtFloat64(t *testing.T) {
	type input struct {
		Value float64
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "zero",
			input: input{
				Value: 0,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": 0.0000000000",
				ResultRemote: "\"name\":0.0000000000",
			},
		},

		{
			desc: "positive whole",
			input: input{
				Value: 1,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": 1.0000000000",
				ResultRemote: "\"name\":1.0000000000",
			},
		},

		{
			desc: "positive with dp",
			input: input{
				Value: 1.1234567890,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": 1.1234567890",
				ResultRemote: "\"name\":1.1234567890",
			},
		},

		{
			desc: "positive with dp greater than 10 round down",
			input: input{
				Value: 1.12345678904,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": 1.1234567890",
				ResultRemote: "\"name\":1.1234567890",
			},
		},

		{
			desc: "positive with dp greater than 10 round up",
			input: input{
				Value: 1.12345678905,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": 1.1234567891",
				ResultRemote: "\"name\":1.1234567891",
			},
		},

		{
			desc: "negative whole",
			input: input{
				Value: -1,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": -1.0000000000",
				ResultRemote: "\"name\":-1.0000000000",
			},
		},

		{
			desc: "negative with dp",
			input: input{
				Value: -1.1234567890,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": -1.1234567890",
				ResultRemote: "\"name\":-1.1234567890",
			},
		},

		{
			desc: "negative with dp greater than 10 round down",
			input: input{
				Value: -1.12345678904,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": -1.1234567890",
				ResultRemote: "\"name\":-1.1234567890",
			},
		},

		{
			desc: "negative with dp greater than 10 round up",
			input: input{
				Value: -1.12345678905,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": -1.1234567891",
				ResultRemote: "\"name\":-1.1234567891",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtFloat64(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtFloat64(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtFloat64s(t *testing.T) {
	type input struct {
		Value []float64
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "single",
			input: input{
				Value: []float64{1.2},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\t1.2000000000\n\t]",
				ResultRemote: "\"name\":[1.2000000000]",
			},
		},

		{
			desc: "multi",
			input: input{
				Value: []float64{1.2, 3.4},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\t1.2000000000,\n\t\t3.4000000000\n\t]",
				ResultRemote: "\"name\":[1.2000000000,3.4000000000]",
			},
		},

		{
			desc: "empty",
			input: input{
				Value: []float64{},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},

		{
			desc: "nil",
			input: input{
				Value: nil,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtFloat64s(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtFloat64s(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtInt(t *testing.T) {
	type input struct {
		Value int
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "zero",
			input: input{
				Value: 0,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": 0",
				ResultRemote: "\"name\":0",
			},
		},

		{
			desc: "positive",
			input: input{
				Value: 1,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": 1",
				ResultRemote: "\"name\":1",
			},
		},

		{
			desc: "negative",
			input: input{
				Value: -1,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": -1",
				ResultRemote: "\"name\":-1",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtInt(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtInt(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtInts(t *testing.T) {
	type input struct {
		Value []int
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "single",
			input: input{
				Value: []int{1},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\t1\n\t]",
				ResultRemote: "\"name\":[1]",
			},
		},

		{
			desc: "multi",
			input: input{
				Value: []int{1, 2},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\t1,\n\t\t2\n\t]",
				ResultRemote: "\"name\":[1,2]",
			},
		},

		{
			desc: "empty",
			input: input{
				Value: []int{},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},

		{
			desc: "nil",
			input: input{
				Value: nil,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtInts(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtInts(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtInt32(t *testing.T) {
	type input struct {
		Value int32
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "zero",
			input: input{
				Value: 0,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": 0",
				ResultRemote: "\"name\":0",
			},
		},

		{
			desc: "positive",
			input: input{
				Value: 1,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": 1",
				ResultRemote: "\"name\":1",
			},
		},

		{
			desc: "negative",
			input: input{
				Value: -1,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": -1",
				ResultRemote: "\"name\":-1",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtInt32(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtInt32(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtInt32s(t *testing.T) {
	type input struct {
		Value []int32
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "single",
			input: input{
				Value: []int32{1},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\t1\n\t]",
				ResultRemote: "\"name\":[1]",
			},
		},

		{
			desc: "multi",
			input: input{
				Value: []int32{1, 2},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\t1,\n\t\t2\n\t]",
				ResultRemote: "\"name\":[1,2]",
			},
		},

		{
			desc: "empty",
			input: input{
				Value: []int32{},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},

		{
			desc: "nil",
			input: input{
				Value: nil,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtInt32s(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtInt32s(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtInt64(t *testing.T) {
	type input struct {
		Value int64
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "zero",
			input: input{
				Value: 0,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": 0",
				ResultRemote: "\"name\":0",
			},
		},

		{
			desc: "positive",
			input: input{
				Value: 1,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": 1",
				ResultRemote: "\"name\":1",
			},
		},

		{
			desc: "negative",
			input: input{
				Value: -1,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": -1",
				ResultRemote: "\"name\":-1",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtInt64(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtInt64(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtInt64s(t *testing.T) {
	type input struct {
		Value []int64
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "single",
			input: input{
				Value: []int64{1},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\t1\n\t]",
				ResultRemote: "\"name\":[1]",
			},
		},

		{
			desc: "multi",
			input: input{
				Value: []int64{1, 2},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\t1,\n\t\t2\n\t]",
				ResultRemote: "\"name\":[1,2]",
			},
		},

		{
			desc: "empty",
			input: input{
				Value: []int64{},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},

		{
			desc: "nil",
			input: input{
				Value: nil,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtInt64s(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtInt64s(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtString(t *testing.T) {
	type input struct {
		Value string
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "string",
			input: input{
				Value: "string",
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": \"string\"",
				ResultRemote: "\"name\":\"string\"",
			},
		},

		{
			desc: "empty",
			input: input{
				Value: "",
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": \"\"",
				ResultRemote: "\"name\":\"\"",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtString(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtString(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtStrings(t *testing.T) {
	type input struct {
		Value []string
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "single",
			input: input{
				Value: []string{"string"},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\t\"string\"\n\t]",
				ResultRemote: "\"name\":[\"string\"]",
			},
		},

		{
			desc: "multi",
			input: input{
				Value: []string{"string", ""},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\t\"string\",\n\t\t\"\"\n\t]",
				ResultRemote: "\"name\":[\"string\",\"\"]",
			},
		},

		{
			desc: "empty",
			input: input{
				Value: []string{},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},

		{
			desc: "nil",
			input: input{
				Value: nil,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtStrings(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtStrings(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtTime(t *testing.T) {
	parsedTime, err := time.Parse(time.RFC3339Nano, "2006-01-02T15:04:05.999999999Z")
	if err != nil {
		t.Fatal("Failed parsing time as RFC3339Nano", err)
	}
	type input struct {
		Value time.Time
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "time",
			input: input{
				Value: parsedTime,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": \"2006-01-02T15:04:05.999999999Z\"",
				ResultRemote: "\"name\":\"2006-01-02T15:04:05.999999999Z\"",
			},
		},

		{
			desc: "zero time",
			input: input{
				Value: time.Time{},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": \"0001-01-01T00:00:00Z\"",
				ResultRemote: "\"name\":\"0001-01-01T00:00:00Z\"",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtTime(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtTime(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtTimes(t *testing.T) {
	parsedTime, err := time.Parse(time.RFC3339Nano, "2006-01-02T15:04:05.999999999Z")
	if err != nil {
		t.Fatal("Failed parsing time as RFC3339Nano", err)
	}
	type input struct {
		Value []time.Time
		Name  string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "single",
			input: input{
				Value: []time.Time{parsedTime},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\t\"2006-01-02T15:04:05.999999999Z\"\n\t]",
				ResultRemote: "\"name\":[\"2006-01-02T15:04:05.999999999Z\"]",
			},
		},

		{
			desc: "multi",
			input: input{
				Value: []time.Time{parsedTime, time.Time{}},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\t\"2006-01-02T15:04:05.999999999Z\",\n\t\t\"0001-01-01T00:00:00Z\"\n\t]",
				ResultRemote: "\"name\":[\"2006-01-02T15:04:05.999999999Z\",\"0001-01-01T00:00:00Z\"]",
			},
		},

		{
			desc: "empty",
			input: input{
				Value: []time.Time{},
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},

		{
			desc: "nil",
			input: input{
				Value: nil,
				Name:  "name",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := FmtTimes(d.input.Value, d.input.Name)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := FmtTimes(d.input.Value, d.input.Name)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_FmtSlice(t *testing.T) {
	type input struct {
		Value  []interface{}
		Name   string
		format string
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "single",
			input: input{
				Value:  []interface{}{"string"},
				Name:   "name",
				format: "%q",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\t\"string\"\n\t]",
				ResultRemote: "\"name\":[\"string\"]",
			},
		},

		{
			desc: "multi",
			input: input{
				Value:  []interface{}{"string", "also string"},
				Name:   "name",
				format: "%q",
			},
			expected: expected{
				Result:       "\"name\": [\n\t\t\"string\",\n\t\t\"also string\"\n\t]",
				ResultRemote: "\"name\":[\"string\",\"also string\"]",
			},
		},

		{
			desc: "empty",
			input: input{
				Value:  []interface{}{},
				Name:   "name",
				format: "",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},

		{
			desc: "nil",
			input: input{
				Value:  nil,
				Name:   "name",
				format: "",
			},
			expected: expected{
				Result:       "\"name\": []",
				ResultRemote: "\"name\":[]",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := fmtSlice(d.input.Value, d.input.Name, d.input.format)

		if string(result) != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     string(result),
			}))
		}

		remote = true
		resultRemote := fmtSlice(d.input.Value, d.input.Name, d.input.format)

		if string(resultRemote) != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(resultRemote)",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     string(resultRemote),
			}))
		}
	}
}

func Test_runtimeLineAndFuncName(t *testing.T) {
	line, funcName := runtimeLineAndFuncName(0)
	pc, _, l, _ := runtime.Caller(0)
	expectedLine := l - 1
	expectedFuncName := runtime.FuncForPC(pc).Name()

	if line != expectedLine {
		t.Error(go_testing.Errorf(go_testing.Error{
			Unexpected: "funcName",
			Expected:   expectedLine,
			Result:     line,
		}))
	}
	if funcName != expectedFuncName {
		t.Error(go_testing.Errorf(go_testing.Error{
			Unexpected: "funcName",
			Expected:   expectedFuncName,
			Result:     funcName,
		}))
	}
}

func Test_FmtFields(t *testing.T) {
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    []Field
		expected expected
	}{
		{
			desc: "single",
			input: []Field{
				Field("field"),
			},
			expected: expected{
				Result:       "{\n\tfield\n}",
				ResultRemote: "{field}",
			},
		},

		{
			desc: "multi",
			input: []Field{
				Field("field"),
				Field("also_field"),
			},
			expected: expected{
				Result:       "{\n\tfield,\n\talso_field\n}",
				ResultRemote: "{field,also_field}",
			},
		},

		{
			desc:  "empty",
			input: []Field{},
			expected: expected{
				Result:       "",
				ResultRemote: "",
			},
		},

		{
			desc:  "nil",
			input: nil,
			expected: expected{
				Result:       "",
				ResultRemote: "",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := fmtFields(d.input)

		if result != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     result,
			}))
		}

		remote = true
		resultRemote := fmtFields(d.input)

		if resultRemote != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "resultRemote",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     resultRemote,
			}))
		}
	}
}

func Test_FmtLog(t *testing.T) {
	type input struct {
		Message       string
		CorrelationId string
		FuncName      string
		Line          int
		Fields        []Field
		Remote        bool
	}
	type expected struct {
		Result       string
		ResultRemote string
	}
	var data = []struct {
		desc     string
		input    input
		expected expected
	}{
		{
			desc: "single",
			input: input{
				Message:       "message",
				CorrelationId: "correlationId",
				FuncName:      "funcName",
				Line:          0,
				Fields: []Field{
					Field("field"),
				},
			},
			expected: expected{
				Result:       "[message] [correlationId] [funcName:0] {\n\tfield\n}\x1b[0m",
				ResultRemote: "[message] [correlationId] [funcName:0] {field}\x1b[0m",
			},
		},

		{
			desc: "multi",
			input: input{
				Message:       "message",
				CorrelationId: "correlationId",
				FuncName:      "funcName",
				Line:          0,
				Fields: []Field{
					Field("field"),
					Field("also_field"),
				},
			},
			expected: expected{
				Result:       "[message] [correlationId] [funcName:0] {\n\tfield,\n\talso_field\n}\x1b[0m",
				ResultRemote: "[message] [correlationId] [funcName:0] {field,also_field}\x1b[0m",
			},
		},

		{
			desc: "empty",
			input: input{
				Message:       "message",
				CorrelationId: "correlationId",
				FuncName:      "funcName",
				Line:          0,
				Fields:        []Field{},
			},
			expected: expected{
				Result:       "[message] [correlationId] [funcName:0] \x1b[0m",
				ResultRemote: "[message] [correlationId] [funcName:0] \x1b[0m",
			},
		},

		{
			desc: "nil",
			input: input{
				Message:       "message",
				CorrelationId: "correlationId",
				FuncName:      "funcName",
				Line:          0,
				Fields:        nil,
			},
			expected: expected{
				Result:       "[message] [correlationId] [funcName:0] \x1b[0m",
				ResultRemote: "[message] [correlationId] [funcName:0] \x1b[0m",
			},
		},
	}

	for i, d := range data {
		remote = false
		result := fmtLog(d.input.Message, d.input.CorrelationId, d.input.FuncName, d.input.Line, d.input.Fields)

		if result != d.expected.Result {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.Result,
				Result:     result,
			}))
		}

		remote = true
		resultRemote := fmtLog(d.input.Message, d.input.CorrelationId, d.input.FuncName, d.input.Line, d.input.Fields)

		if resultRemote != d.expected.ResultRemote {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "resultRemote",
				Desc:       d.desc,
				At:         i,
				Input:      d.input,
				Expected:   d.expected.ResultRemote,
				Result:     resultRemote,
			}))
		}
	}
}
