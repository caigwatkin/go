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
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"

	go_errors "github.com/caigwatkin/go/errors"
	go_testing "github.com/caigwatkin/go/testing"
)

func TestNewClient(t *testing.T) {
	var data = []struct {
		desc     string
		input    bool
		expected client
	}{
		{
			desc:  "debug enabled",
			input: true,
			expected: client{
				debug: true,
			},
		},

		{
			desc:  "debug disabled",
			input: false,
			expected: client{
				debug: false,
			},
		},
	}

	for i, d := range data {
		result := NewClient(d.input)

		if reflect.TypeOf(result) != reflect.TypeOf(d.expected) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result",
				Desc:       d.desc,
				At:         i,
				Expected:   reflect.TypeOf(d.expected),
				Result:     reflect.TypeOf(result),
			}))
		}
		if v, ok := result.(client); !ok {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     result,
			}))

		} else if v.debug != d.expected.debug {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.debug",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.debug,
				Result:     v.debug,
			}))
		}
	}
}

func TestFmtAny(t *testing.T) {
	notJSONMarshallableFunc := func() {}
	type input struct {
		value interface{}
		name  string
	}
	type valueStruct struct {
		X string `json:"x,omitempty"`
		Y int    `json:"y,omitempty"`
		Z bool   `json:"z,omitempty"`
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "struct as indented JSON",
			input: input{
				value: valueStruct{
					X: "some string",
					Y: 2,
					Z: true,
				},
				name: "name",
			},
			expected: "\"name\": {\n\t\t\"type\": \"log.valueStruct\",\n\t\t\"value\": {\n\t\t\t\"x\": \"some string\",\n\t\t\t\"y\": 2,\n\t\t\t\"z\": true\n\t\t}\n\t}",
		},

		{
			desc: "struct omitempty",
			input: input{
				value: valueStruct{
					X: "",
					Y: 0,
					Z: false,
				},
				name: "name",
			},
			expected: "\"name\": {\n\t\t\"type\": \"log.valueStruct\",\n\t\t\"value\": {}\n\t}",
		},

		{
			desc: "not JSON marshallable",
			input: input{
				value: notJSONMarshallableFunc,
				name:  "name",
			},
			expected: "\"name\": {\n\t\t\"type\": \"func()\",\n\t\t\"value\": \"NOT JSON MARSHALLABLE\"\n\t}",
		},

		{
			desc: "nil",
			input: input{
				value: nil,
				name:  "name",
			},
			expected: "\"name\": null",
		},
	}

	for i, d := range data {
		result := FmtAny(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtAnys(t *testing.T) {
	notJSONMarshallableFunc := func() {}
	type input struct {
		value []interface{}
		name  string
	}
	type valueStruct struct {
		X string `json:"x,omitempty"`
		Y int    `json:"y,omitempty"`
		Z bool   `json:"z,omitempty"`
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "single",
			input: input{
				value: []interface{}{
					valueStruct{
						X: "some string",
						Y: 2,
						Z: true,
					},
				},
				name: "name",
			},
			expected: "\"name\": [\n\t\t{\n\t\t\t\"type\": \"log.valueStruct\",\n\t\t\t\"value\": {\n\t\t\t\t\"x\": \"some string\",\n\t\t\t\t\"y\": 2,\n\t\t\t\t\"z\": true\n\t\t\t}\n\t\t}\n\t]",
		},

		{
			desc: "multi",
			input: input{
				value: []interface{}{
					valueStruct{
						X: "some string",
						Y: 2,
						Z: true,
					},
					nil,
					notJSONMarshallableFunc,
				},
				name: "name",
			},
			expected: "\"name\": [\n\t\t{\n\t\t\t\"type\": \"log.valueStruct\",\n\t\t\t\"value\": {\n\t\t\t\t\"x\": \"some string\",\n\t\t\t\t\"y\": 2,\n\t\t\t\t\"z\": true\n\t\t\t}\n\t\t},\n\t\tnull,\n\t\t{\n\t\t\t\"type\": \"func()\",\n\t\t\t\"value\": \"NOT JSON MARSHALLABLE\"\n\t\t}\n\t]",
		},

		{
			desc: "empty",
			input: input{
				value: []interface{}{},
				name:  "name",
			},
			expected: "\"name\": []",
		},

		{
			desc: "nil",
			input: input{
				value: nil,
				name:  "name",
			},
			expected: "\"name\": []",
		},
	}

	for i, d := range data {
		result := FmtAnys(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtBool(t *testing.T) {
	type input struct {
		value bool
		name  string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "false",
			input: input{
				value: false,
				name:  "name",
			},
			expected: "\"name\": false",
		},

		{
			desc: "true",
			input: input{
				value: true,
				name:  "name",
			},
			expected: "\"name\": true",
		},
	}

	for i, d := range data {
		result := FmtBool(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtBools(t *testing.T) {
	type input struct {
		value []bool
		name  string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "single",
			input: input{
				value: []bool{true},
				name:  "name",
			},
			expected: "\"name\": [\n\t\ttrue\n\t]",
		},

		{
			desc: "multi",
			input: input{
				value: []bool{true, false},
				name:  "name",
			},
			expected: "\"name\": [\n\t\ttrue,\n\t\tfalse\n\t]",
		},

		{
			desc: "empty",
			input: input{
				value: []bool{},
				name:  "name",
			},
			expected: "\"name\": []",
		},

		{
			desc: "nil",
			input: input{
				value: nil,
				name:  "name",
			},
			expected: "\"name\": []",
		},
	}

	for i, d := range data {
		result := FmtBools(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtByte(t *testing.T) {
	type input struct {
		value byte
		name  string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "Character",
			input: input{
				value: 'A',
				name:  "name",
			},
			expected: "\"name\": 'A'",
		},

		{
			desc: "integer",
			input: input{
				value: 0,
				name:  "name",
			},
			expected: "\"name\": '\\x00'",
		},
	}

	for i, d := range data {
		result := FmtByte(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtBytes(t *testing.T) {
	type input struct {
		value []byte
		name  string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "nil",
			input: input{
				value: nil,
				name:  "name",
			},
			expected: "\"name\": \"\"",
		},

		{
			desc: "empty",
			input: input{
				value: []byte{},
				name:  "name",
			},
			expected: "\"name\": \"\"",
		},

		{
			desc: "string",
			input: input{
				value: []byte("some text"),
				name:  "name",
			},
			expected: "\"name\": \"some text\"",
		},

		{
			desc: "characters",
			input: input{
				value: []byte{'A', 'a', 'B'},
				name:  "name",
			},
			expected: "\"name\": \"AaB\"",
		},

		{
			desc: "characters",
			input: input{
				value: []byte{0, 25, 18},
				name:  "name",
			},
			expected: "\"name\": \"\\x00\\x19\\x12\"",
		},
	}

	for i, d := range data {
		result := FmtBytes(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtDuration(t *testing.T) {
	type input struct {
		value time.Duration
		name  string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "nano",
			input: input{
				value: time.Nanosecond,
				name:  "name",
			},
			expected: "\"name\": \"1ns\"",
		},

		{
			desc: "micro",
			input: input{
				value: time.Microsecond,
				name:  "name",
			},
			expected: "\"name\": \"1µs\"",
		},

		{
			desc: "milli",
			input: input{
				value: time.Millisecond,
				name:  "name",
			},
			expected: "\"name\": \"1ms\"",
		},

		{
			desc: "second",
			input: input{
				value: time.Second,
				name:  "name",
			},
			expected: "\"name\": \"1s\"",
		},

		{
			desc: "minute",
			input: input{
				value: time.Minute,
				name:  "name",
			},
			expected: "\"name\": \"1m0s\"",
		},

		{
			desc: "hour",
			input: input{
				value: time.Hour,
				name:  "name",
			},
			expected: "\"name\": \"1h0m0s\"",
		},
	}

	for i, d := range data {
		result := FmtDuration(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtDurations(t *testing.T) {
	type input struct {
		value []time.Duration
		name  string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "single",
			input: input{
				value: []time.Duration{time.Nanosecond},
				name:  "name",
			},
			expected: "\"name\": [\n\t\t\"1ns\"\n\t]",
		},

		{
			desc: "multi",
			input: input{
				value: []time.Duration{time.Nanosecond, time.Microsecond},
				name:  "name",
			},
			expected: "\"name\": [\n\t\t\"1ns\",\n\t\t\"1µs\"\n\t]",
		},

		{
			desc: "empty",
			input: input{
				value: []time.Duration{},
				name:  "name",
			},
			expected: "\"name\": []",
		},

		{
			desc: "nil",
			input: input{
				value: nil,
				name:  "name",
			},
			expected: "\"name\": []",
		},
	}

	for i, d := range data {
		result := FmtDurations(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtError(t *testing.T) {
	errWithTrace := go_errors.New("error")
	trace := fmt.Sprintf("%+v", errWithTrace)
	var data = []struct {
		desc     string
		input    error
		expected string
	}{
		{
			desc:     "nil",
			input:    nil,
			expected: "\"error\": null",
		},

		{
			desc:     "no trace",
			input:    errors.New("some error"),
			expected: "\"error\": \"some error\"",
		},

		{
			desc:     "trace",
			input:    errWithTrace,
			expected: fmt.Sprintf("\"error\": {\n\t\t\"friendly\": \"error\",\n\t\t\"trace\": %s\n\t}", trace),
		},
	}

	for i, d := range data {
		result := FmtError(d.input)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtFloat32(t *testing.T) {
	type input struct {
		value float32
		name  string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "zero",
			input: input{
				value: 0,
				name:  "name",
			},
			expected: "\"name\": 0.00000",
		},

		{
			desc: "positive whole",
			input: input{
				value: 1,
				name:  "name",
			},
			expected: "\"name\": 1.00000",
		},

		{
			desc: "positive with dp",
			input: input{
				value: 1.12345,
				name:  "name",
			},
			expected: "\"name\": 1.12345",
		},

		{
			desc: "positive with dp greater than 5 round down",
			input: input{
				value: 1.123454,
				name:  "name",
			},
			expected: "\"name\": 1.12345",
		},

		{
			desc: "positive with dp greater than 5 round up",
			input: input{
				value: 1.123455,
				name:  "name",
			},
			expected: "\"name\": 1.12346",
		},

		{
			desc: "negative whole",
			input: input{
				value: -1,
				name:  "name",
			},
			expected: "\"name\": -1.00000",
		},

		{
			desc: "negative with dp",
			input: input{
				value: -1.12345,
				name:  "name",
			},
			expected: "\"name\": -1.12345",
		},

		{
			desc: "negative with dp greater than 5 round down",
			input: input{
				value: -1.123454,
				name:  "name",
			},
			expected: "\"name\": -1.12345",
		},

		{
			desc: "negative with dp greater than 5 round up",
			input: input{
				value: -1.123455,
				name:  "name",
			},
			expected: "\"name\": -1.12346",
		},
	}

	for i, d := range data {
		result := FmtFloat32(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtFloat32s(t *testing.T) {
	type input struct {
		value []float32
		name  string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "single",
			input: input{
				value: []float32{1.2},
				name:  "name",
			},
			expected: "\"name\": [\n\t\t1.20000\n\t]",
		},

		{
			desc: "multi",
			input: input{
				value: []float32{1.2, 3.4},
				name:  "name",
			},
			expected: "\"name\": [\n\t\t1.20000,\n\t\t3.40000\n\t]",
		},

		{
			desc: "empty",
			input: input{
				value: []float32{},
				name:  "name",
			},
			expected: "\"name\": []",
		},

		{
			desc: "nil",
			input: input{
				value: nil,
				name:  "name",
			},
			expected: "\"name\": []",
		},
	}

	for i, d := range data {
		result := FmtFloat32s(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtFloat64(t *testing.T) {
	type input struct {
		value float64
		name  string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "zero",
			input: input{
				value: 0,
				name:  "name",
			},
			expected: "\"name\": 0.0000000000",
		},

		{
			desc: "positive whole",
			input: input{
				value: 1,
				name:  "name",
			},
			expected: "\"name\": 1.0000000000",
		},

		{
			desc: "positive with dp",
			input: input{
				value: 1.1234567890,
				name:  "name",
			},
			expected: "\"name\": 1.1234567890",
		},

		{
			desc: "positive with dp greater than 10 round down",
			input: input{
				value: 1.12345678904,
				name:  "name",
			},
			expected: "\"name\": 1.1234567890",
		},

		{
			desc: "positive with dp greater than 10 round up",
			input: input{
				value: 1.12345678905,
				name:  "name",
			},
			expected: "\"name\": 1.1234567891",
		},

		{
			desc: "negative whole",
			input: input{
				value: -1,
				name:  "name",
			},
			expected: "\"name\": -1.0000000000",
		},

		{
			desc: "negative with dp",
			input: input{
				value: -1.1234567890,
				name:  "name",
			},
			expected: "\"name\": -1.1234567890",
		},

		{
			desc: "negative with dp greater than 10 round down",
			input: input{
				value: -1.12345678904,
				name:  "name",
			},
			expected: "\"name\": -1.1234567890",
		},

		{
			desc: "negative with dp greater than 10 round up",
			input: input{
				value: -1.12345678905,
				name:  "name",
			},
			expected: "\"name\": -1.1234567891",
		},
	}

	for i, d := range data {
		result := FmtFloat64(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtFloat64s(t *testing.T) {
	type input struct {
		value []float64
		name  string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "single",
			input: input{
				value: []float64{1.2},
				name:  "name",
			},
			expected: "\"name\": [\n\t\t1.2000000000\n\t]",
		},

		{
			desc: "multi",
			input: input{
				value: []float64{1.2, 3.4},
				name:  "name",
			},
			expected: "\"name\": [\n\t\t1.2000000000,\n\t\t3.4000000000\n\t]",
		},

		{
			desc: "empty",
			input: input{
				value: []float64{},
				name:  "name",
			},
			expected: "\"name\": []",
		},

		{
			desc: "nil",
			input: input{
				value: nil,
				name:  "name",
			},
			expected: "\"name\": []",
		},
	}

	for i, d := range data {
		result := FmtFloat64s(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtInt(t *testing.T) {
	type input struct {
		value int
		name  string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "zero",
			input: input{
				value: 0,
				name:  "name",
			},
			expected: "\"name\": 0",
		},

		{
			desc: "positive",
			input: input{
				value: 1,
				name:  "name",
			},
			expected: "\"name\": 1",
		},

		{
			desc: "negative",
			input: input{
				value: -1,
				name:  "name",
			},
			expected: "\"name\": -1",
		},
	}

	for i, d := range data {
		result := FmtInt(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtInts(t *testing.T) {
	type input struct {
		value []int
		name  string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "single",
			input: input{
				value: []int{1},
				name:  "name",
			},
			expected: "\"name\": [\n\t\t1\n\t]",
		},

		{
			desc: "multi",
			input: input{
				value: []int{1, 2},
				name:  "name",
			},
			expected: "\"name\": [\n\t\t1,\n\t\t2\n\t]",
		},

		{
			desc: "empty",
			input: input{
				value: []int{},
				name:  "name",
			},
			expected: "\"name\": []",
		},

		{
			desc: "nil",
			input: input{
				value: nil,
				name:  "name",
			},
			expected: "\"name\": []",
		},
	}

	for i, d := range data {
		result := FmtInts(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtInt32(t *testing.T) {
	type input struct {
		value int32
		name  string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "zero",
			input: input{
				value: 0,
				name:  "name",
			},
			expected: "\"name\": 0",
		},

		{
			desc: "positive",
			input: input{
				value: 1,
				name:  "name",
			},
			expected: "\"name\": 1",
		},

		{
			desc: "negative",
			input: input{
				value: -1,
				name:  "name",
			},
			expected: "\"name\": -1",
		},
	}

	for i, d := range data {
		result := FmtInt32(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtInt32s(t *testing.T) {
	type input struct {
		value []int32
		name  string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "single",
			input: input{
				value: []int32{1},
				name:  "name",
			},
			expected: "\"name\": [\n\t\t1\n\t]",
		},

		{
			desc: "multi",
			input: input{
				value: []int32{1, 2},
				name:  "name",
			},
			expected: "\"name\": [\n\t\t1,\n\t\t2\n\t]",
		},

		{
			desc: "empty",
			input: input{
				value: []int32{},
				name:  "name",
			},
			expected: "\"name\": []",
		},

		{
			desc: "nil",
			input: input{
				value: nil,
				name:  "name",
			},
			expected: "\"name\": []",
		},
	}

	for i, d := range data {
		result := FmtInt32s(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtInt64(t *testing.T) {
	type input struct {
		value int64
		name  string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "zero",
			input: input{
				value: 0,
				name:  "name",
			},
			expected: "\"name\": 0",
		},

		{
			desc: "positive",
			input: input{
				value: 1,
				name:  "name",
			},
			expected: "\"name\": 1",
		},

		{
			desc: "negative",
			input: input{
				value: -1,
				name:  "name",
			},
			expected: "\"name\": -1",
		},
	}

	for i, d := range data {
		result := FmtInt64(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtInt64s(t *testing.T) {
	type input struct {
		value []int64
		name  string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "single",
			input: input{
				value: []int64{1},
				name:  "name",
			},
			expected: "\"name\": [\n\t\t1\n\t]",
		},

		{
			desc: "multi",
			input: input{
				value: []int64{1, 2},
				name:  "name",
			},
			expected: "\"name\": [\n\t\t1,\n\t\t2\n\t]",
		},

		{
			desc: "empty",
			input: input{
				value: []int64{},
				name:  "name",
			},
			expected: "\"name\": []",
		},

		{
			desc: "nil",
			input: input{
				value: nil,
				name:  "name",
			},
			expected: "\"name\": []",
		},
	}

	for i, d := range data {
		result := FmtInt64s(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtString(t *testing.T) {
	type input struct {
		value string
		name  string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "string",
			input: input{
				value: "string",
				name:  "name",
			},
			expected: "\"name\": \"string\"",
		},

		{
			desc: "empty",
			input: input{
				value: "",
				name:  "name",
			},
			expected: "\"name\": \"\"",
		},
	}

	for i, d := range data {
		result := FmtString(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtStrings(t *testing.T) {
	type input struct {
		value []string
		name  string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "single",
			input: input{
				value: []string{"string"},
				name:  "name",
			},
			expected: "\"name\": [\n\t\t\"string\"\n\t]",
		},

		{
			desc: "multi",
			input: input{
				value: []string{"string", ""},
				name:  "name",
			},
			expected: "\"name\": [\n\t\t\"string\",\n\t\t\"\"\n\t]",
		},

		{
			desc: "empty",
			input: input{
				value: []string{},
				name:  "name",
			},
			expected: "\"name\": []",
		},

		{
			desc: "nil",
			input: input{
				value: nil,
				name:  "name",
			},
			expected: "\"name\": []",
		},
	}

	for i, d := range data {
		result := FmtStrings(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtTime(t *testing.T) {
	parsedTime, err := time.Parse(time.RFC3339Nano, "2006-01-02T15:04:05.999999999Z")
	if err != nil {
		t.Fatal("Failed parsing time as RFC3339Nano", err)
	}
	type input struct {
		value time.Time
		name  string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "time",
			input: input{
				value: parsedTime,
				name:  "name",
			},
			expected: "\"name\": \"2006-01-02T15:04:05.999999999Z\"",
		},

		{
			desc: "zero time",
			input: input{
				value: time.Time{},
				name:  "name",
			},
			expected: "\"name\": \"0001-01-01T00:00:00Z\"",
		},
	}

	for i, d := range data {
		result := FmtTime(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtTimes(t *testing.T) {
	parsedTime, err := time.Parse(time.RFC3339Nano, "2006-01-02T15:04:05.999999999Z")
	if err != nil {
		t.Fatal("Failed parsing time as RFC3339Nano", err)
	}
	type input struct {
		value []time.Time
		name  string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "single",
			input: input{
				value: []time.Time{parsedTime},
				name:  "name",
			},
			expected: "\"name\": [\n\t\t\"2006-01-02T15:04:05.999999999Z\"\n\t]",
		},

		{
			desc: "multi",
			input: input{
				value: []time.Time{parsedTime, time.Time{}},
				name:  "name",
			},
			expected: "\"name\": [\n\t\t\"2006-01-02T15:04:05.999999999Z\",\n\t\t\"0001-01-01T00:00:00Z\"\n\t]",
		},

		{
			desc: "empty",
			input: input{
				value: []time.Time{},
				name:  "name",
			},
			expected: "\"name\": []",
		},

		{
			desc: "nil",
			input: input{
				value: nil,
				name:  "name",
			},
			expected: "\"name\": []",
		},
	}

	for i, d := range data {
		result := FmtTimes(d.input.value, d.input.name)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestFmtSlice(t *testing.T) {
	type input struct {
		value  []interface{}
		name   string
		format string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "single",
			input: input{
				value:  []interface{}{"string"},
				name:   "name",
				format: "%q",
			},
			expected: "\"name\": [\n\t\t\"string\"\n\t]",
		},

		{
			desc: "multi",
			input: input{
				value:  []interface{}{"string", "also string"},
				name:   "name",
				format: "%q",
			},
			expected: "\"name\": [\n\t\t\"string\",\n\t\t\"also string\"\n\t]",
		},

		{
			desc: "empty",
			input: input{
				value:  []interface{}{},
				name:   "name",
				format: "",
			},
			expected: "\"name\": []",
		},

		{
			desc: "nil",
			input: input{
				value:  nil,
				name:   "name",
				format: "",
			},
			expected: "\"name\": []",
		},
	}

	for i, d := range data {
		result := fmtSlice(d.input.value, d.input.name, d.input.format)

		if string(result) != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "string(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}

func TestRuntimeLineAndFuncName(t *testing.T) {
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

func TestFmtFields(t *testing.T) {
	var data = []struct {
		desc     string
		input    []Field
		expected string
	}{
		{
			desc: "single",
			input: []Field{
				Field("field"),
			},
			expected: "{\n\tfield\n}",
		},

		{
			desc: "multi",
			input: []Field{
				Field("field"),
				Field("also field"),
			},
			expected: "{\n\tfield,\n\talso field\n}",
		},

		{
			desc:     "empty",
			input:    []Field{},
			expected: "",
		},

		{
			desc:     "nil",
			input:    nil,
			expected: "",
		},
	}

	for i, d := range data {
		result := fmtFields(d.input)

		if result != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     result,
			}))
		}
	}
}

func TestFmtLog(t *testing.T) {
	type input struct {
		message       string
		correlationID string
		funcName      string
		line          int
		fields        []Field
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "single",
			input: input{
				message:       "message",
				correlationID: "correlationID",
				funcName:      "funcName",
				line:          0,
				fields: []Field{
					Field("field"),
				},
			},
			expected: "[message] [correlationID] [funcName:0] {\n\tfield\n}\x1b[0m",
		},

		{
			desc: "multi",
			input: input{
				message:       "message",
				correlationID: "correlationID",
				funcName:      "funcName",
				line:          0,
				fields: []Field{
					Field("field"),
					Field("also field"),
				},
			},
			expected: "[message] [correlationID] [funcName:0] {\n\tfield,\n\talso field\n}\x1b[0m",
		},

		{
			desc: "empty",
			input: input{
				message:       "message",
				correlationID: "correlationID",
				funcName:      "funcName",
				line:          0,
				fields:        []Field{},
			},
			expected: "[message] [correlationID] [funcName:0] \x1b[0m",
		},

		{
			desc: "nil",
			input: input{
				message:       "message",
				correlationID: "correlationID",
				funcName:      "funcName",
				line:          0,
				fields:        nil,
			},
			expected: "[message] [correlationID] [funcName:0] \x1b[0m",
		},
	}

	for i, d := range data {
		result := fmtLog(d.input.message, d.input.correlationID, d.input.funcName, d.input.line, d.input.fields)

		if result != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     result,
			}))
		}
	}
}
