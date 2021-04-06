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

		if !reflect.DeepEqual(result.Code, d.expected.Code) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Code",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.Code,
				Result:     result.Code,
			}))
		}
		if !reflect.DeepEqual(result.Message, d.expected.Message) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Message",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.Message,
				Result:     result.Message,
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

		if !reflect.DeepEqual(result.Code, d.expected.Code) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Code",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.Code,
				Result:     result.Code,
			}))
		}
		if !reflect.DeepEqual(result.Message, d.expected.Message) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Message",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.Message,
				Result:     result.Message,
			}))
		}
	}
}

func Test_NewStatusWithCause(t *testing.T) {
	err := New("error")
	type input struct {
		Code    int
		Message string
		Cause   error
	}
	var data = []struct {
		desc string
		input
		expected Status
	}{
		{
			desc: "status no cause",
			input: input{
				Code:    http.StatusAccepted,
				Message: "This has been accepted",
				Cause:   nil,
			},
			expected: Status{
				Code:    http.StatusAccepted,
				Message: "Accepted: This has been accepted",
				Cause:   nil,
			},
		},

		{
			desc: "status with cause",
			input: input{
				Cause:   err,
				Code:    http.StatusAccepted,
				Message: "This has been accepted",
			},
			expected: Status{
				Cause:   err,
				Code:    http.StatusAccepted,
				Message: "Accepted: This has been accepted",
			},
		},
	}

	for i, d := range data {
		result := NewStatusWithCause(d.input.Cause, d.input.Code, d.input.Message)

		if d.expected.Cause != nil {
			if result.Cause == nil || result.Cause.Error() != d.expected.Cause.Error() {
				t.Error(go_testing.Errorf(go_testing.Error{
					Unexpected: "result.Cause",
					Desc:       d.desc,
					At:         i,
					Expected:   d.expected.Cause,
					Result:     result.Cause,
				}))
			} else if result.Cause.Error() != d.expected.Cause.Error() {
				t.Error(go_testing.Errorf(go_testing.Error{
					Unexpected: "result.Cause.Error()",
					Desc:       d.desc,
					At:         i,
					Expected:   d.expected.Cause.Error(),
					Result:     result.Cause.Error(),
				}))
			}
		} else if result.Cause != nil {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Cause",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.Cause,
				Result:     result.Cause,
			}))
		}
		if !reflect.DeepEqual(result.Code, d.expected.Code) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Code",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.Code,
				Result:     result.Code,
			}))
		}
		if !reflect.DeepEqual(result.Message, d.expected.Message) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Message",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.Message,
				Result:     result.Message,
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
			desc: "status no items",
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
			desc: "status with items",
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

		if !reflect.DeepEqual(result.Code, d.expected.Code) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Code",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.Code,
				Result:     result.Code,
			}))
		}
		if !reflect.DeepEqual(result.Message, d.expected.Message) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Message",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.Message,
				Result:     result.Message,
			}))
		}
		if !reflect.DeepEqual(result.Items, d.expected.Items) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Items",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.Items,
				Result:     result.Items,
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

func Test_Error(t *testing.T) {
	var data = []struct {
		desc     string
		input    Status
		expected string
	}{
		{
			desc:     "defaults",
			input:    Status{},
			expected: "\"code\": 0, \"message\": \"\", \"at\": \"\", \"items\": []",
		},
		{
			desc: "values",
			input: Status{
				At:      "at",
				Code:    http.StatusAccepted,
				Message: "message",
				Items: []Item{
					{
						Field:   "item_field",
						Message: "item_message",
					},
				},
			},
			expected: "\"code\": 202, \"message\": \"message\", \"at\": \"at\", \"items\": [{item_field item_message}]",
		},
	}

	for i, d := range data {
		result := d.input.Error()

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
			desc: "items",
			input: Status{
				Items: []Item{
					{
						Field:   "field",
						Message: "message",
					},
					{
						Field:   "",
						Message: "message_2",
					},
				},
			},
			expected: []byte("[{\"field\":\"field\",\"message\":\"message\"},{\"message\":\"message_2\"}]"),
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
