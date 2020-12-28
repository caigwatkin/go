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

package errors

import (
	"net/http"
	"reflect"
	"testing"

	go_testing "github.com/caigwatkin/go/testing"
)

func Test_NewStatus(t *testing.T) {
	type input struct {
		Code    int
		Message string
	}
	var data = []struct {
		desc string
		input
		expected Status
	}{
		{
			desc: "status",
			input: input{
				Code:    http.StatusAccepted,
				Message: "This has been accepted",
			},
			expected: Status{
				Code:    http.StatusAccepted,
				Message: "Accepted: This has been accepted",
			},
		},
	}

	for i, d := range data {
		result := NewStatus(d.input.Code, d.input.Message)

		if !reflect.DeepEqual(result, d.expected) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Code",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.Code,
				Result:     result.Code,
			}))
		}
	}
}

func Test_Statusf(t *testing.T) {
	type input struct {
		Code    int
		Message string
		Args    []interface{}
	}
	var data = []struct {
		desc string
		input
		expected Status
	}{
		{
			desc: "status no args",
			input: input{
				Code:    http.StatusAccepted,
				Message: "This has been accepted",
				Args:    nil,
			},
			expected: Status{
				Code:    http.StatusAccepted,
				Message: "Accepted: This has been accepted",
			},
		},

		{
			desc: "status with args",
			input: input{
				Code:    http.StatusAccepted,
				Message: "This has been %s",
				Args:    []interface{}{"accepted"},
			},
			expected: Status{
				Code:    http.StatusAccepted,
				Message: "Accepted: This has been accepted",
			},
		},
	}

	for i, d := range data {
		result := Statusf(d.input.Code, d.input.Message, d.input.Args...)

		if !reflect.DeepEqual(result, d.expected) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Code",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.Code,
				Result:     result.Code,
			}))
		}
	}
}

func Test_NewStatusWithItems(t *testing.T) {
	type input struct {
		Code    int
		Message string
		Items   []Item
	}
	var data = []struct {
		desc string
		input
		expected Status
	}{
		{
			desc: "status no metadata",
			input: input{
				Code:    http.StatusAccepted,
				Message: "This has been accepted",
				Items:   nil,
			},
			expected: Status{
				Code:    http.StatusAccepted,
				Message: "Accepted: This has been accepted",
				Items:   nil,
			},
		},

		{
			desc: "status with metadata",
			input: input{
				Code:    http.StatusAccepted,
				Message: "This has been accepted",
				Items: []Item{
					{
						Field:   "field",
						Message: "message",
					},
				},
			},
			expected: Status{
				Code:    http.StatusAccepted,
				Message: "Accepted: This has been accepted",
				Items: []Item{
					{
						Field:   "field",
						Message: "message",
					},
				},
			},
		},
	}

	for i, d := range data {
		result := NewStatusWithItems(d.input.Code, d.input.Message, d.input.Items)

		if !reflect.DeepEqual(result, d.expected) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Code",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.Code,
				Result:     result.Code,
			}))
		}
	}
}

func Test_StatusCode(t *testing.T) {
	var data = []struct {
		desc     string
		input    error
		expected int
	}{
		{
			desc:     "status",
			input:    Status{Code: http.StatusAccepted},
			expected: http.StatusAccepted,
		},
		{

			desc:     "nil",
			input:    nil,
			expected: 0,
		},

		{
			desc:     "error",
			input:    New(""),
			expected: 0,
		},
	}

	for i, d := range data {
		result := StatusCode(d.input)

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

func Test_IsStatus(t *testing.T) {
	var data = []struct {
		desc     string
		input    error
		expected bool
	}{
		{
			desc:     "status",
			input:    Status{},
			expected: true,
		},

		{
			desc:     "nil",
			input:    nil,
			expected: false,
		},
	}

	for i, d := range data {
		result := IsStatus(d.input)

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

func Test_Status_RenderItems(t *testing.T) {
	var data = []struct {
		desc     string
		input    Status
		expected []byte
	}{
		{
			desc: "none",
			input: Status{
				Items: nil,
			},
			expected: nil,
		},
		{
			desc: "empty",
			input: Status{
				Items: []Item{},
			},
			expected: nil,
		},
		{
			desc: "item",
			input: Status{
				Items: []Item{
					{
						Field:   "field",
						Message: "message",
					},
				},
			},
			expected: []byte("[{\"field\":\"field\",\"message\":\"message\"}]"),
		},
		{
			desc: "item",
			input: Status{
				Items: []Item{
					{
						Field:   "field",
						Message: "message",
					},
					{
						Field:   "field_2",
						Message: "message_2",
					},
				},
			},
			expected: []byte("[{\"field\":\"field\",\"message\":\"message\"},{\"field\":\"field_2\",\"message\":\"message_2\"}]"),
		},
	}

	for i, d := range data {
		result := d.input.RenderItems()

		if !reflect.DeepEqual(result, d.expected) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     string(result),
			}))
		}
	}
}
