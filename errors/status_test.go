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
	"testing"

	go_testing "github.com/caigwatkin/go/testing"
)

func TestNewStatus(t *testing.T) {
	type input struct {
		code    int
		message string
	}
	var data = []struct {
		desc string
		input
		expected Status
	}{
		{
			desc: "status",
			input: input{
				code:    http.StatusAccepted,
				message: "This has been accepted",
			},
			expected: Status{
				Code:    http.StatusAccepted,
				Message: "Accepted: This has been accepted",
			},
		},
	}

	for i, d := range data {
		result := NewStatus(d.input.code, d.input.message)

		if result.Code != d.expected.Code {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Code",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.Code,
				Result:     result.Code,
			}))
		}
		if result.Message != d.expected.Message {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Message",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.Message,
				Result:     result.Message,
			}))
		}
		if len(result.Metadata) != len(d.expected.Metadata) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Metadata",
				Desc:       d.desc,
				At:         i,
				Expected:   len(d.expected.Metadata),
				Result:     len(result.Metadata),
			}))
		}
	}
}

func TestStatusf(t *testing.T) {
	type input struct {
		code    int
		message string
		args    []interface{}
	}
	var data = []struct {
		desc string
		input
		expected Status
	}{
		{
			desc: "status no args",
			input: input{
				code:    http.StatusAccepted,
				message: "This has been accepted",
				args:    nil,
			},
			expected: Status{
				Code:    http.StatusAccepted,
				Message: "Accepted: This has been accepted",
			},
		},

		{
			desc: "status with args",
			input: input{
				code:    http.StatusAccepted,
				message: "This has been %s",
				args:    []interface{}{"accepted"},
			},
			expected: Status{
				Code:    http.StatusAccepted,
				Message: "Accepted: This has been accepted",
			},
		},
	}

	for i, d := range data {
		result := Statusf(d.input.code, d.input.message, d.input.args...)

		if result.Code != d.expected.Code {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Code",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.Code,
				Result:     result.Code,
			}))
		}
		if result.Message != d.expected.Message {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Message",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.Message,
				Result:     result.Message,
			}))
		}
		if len(result.Metadata) != len(d.expected.Metadata) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Metadata",
				Desc:       d.desc,
				At:         i,
				Expected:   len(d.expected.Metadata),
				Result:     len(result.Metadata),
			}))
		}
	}
}

func TestNewStatusWithMetadata(t *testing.T) {
	type input struct {
		code     int
		message  string
		metadata map[string]interface{}
	}
	var data = []struct {
		desc string
		input
		expected Status
	}{
		{
			desc: "status no metadata",
			input: input{
				code:     http.StatusAccepted,
				message:  "This has been accepted",
				metadata: nil,
			},
			expected: Status{
				Code:     http.StatusAccepted,
				Message:  "Accepted: This has been accepted",
				Metadata: nil,
			},
		},

		{
			desc: "status with metadata",
			input: input{
				code:     http.StatusAccepted,
				message:  "This has been accepted",
				metadata: map[string]interface{}{"some": "metadata"},
			},
			expected: Status{
				Code:     http.StatusAccepted,
				Message:  "Accepted: This has been accepted",
				Metadata: map[string]interface{}{"some": "metadata"},
			},
		},
	}

	for i, d := range data {
		result := NewStatusWithMetadata(d.input.code, d.input.message, d.input.metadata)

		if result.Code != d.expected.Code {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Code",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.Code,
				Result:     result.Code,
			}))
		}
		if result.Message != d.expected.Message {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Message",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.Message,
				Result:     result.Message,
			}))
		}
		if len(result.Metadata) != len(d.expected.Metadata) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Metadata",
				Desc:       d.desc,
				At:         i,
				Expected:   len(d.expected.Metadata),
				Result:     len(result.Metadata),
			}))
		}
	}
}

func TestStatusCode(t *testing.T) {
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

func TestIsStatus(t *testing.T) {
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

func TestRender(t *testing.T) {
	var data = []struct {
		desc     string
		input    Status
		expected string
	}{
		{
			desc: "status bad request",
			input: Status{
				Code:     http.StatusBadRequest,
				Message:  "Some message",
				Metadata: nil,
			},
			expected: `{
	"code": 400,
	"message": "Some message"
}`,
		},

		{
			desc: "status unprocessable entity",
			input: Status{
				Code:     http.StatusUnprocessableEntity,
				Message:  "Some message",
				Metadata: nil,
			},
			expected: `{
	"code": 422,
	"message": "Some message"
}`,
		},

		{
			desc: "status other, message overwritten with status text",
			input: Status{
				Code:     http.StatusAccepted,
				Message:  "Some message",
				Metadata: nil,
			},
			expected: `{
	"code": 202,
	"message": "Accepted"
}`,
		},

		{
			desc: "status 0, no message",
			input: Status{
				Code:     0,
				Message:  "",
				Metadata: nil,
			},
			expected: `{
	"code": 0,
	"message": ""
}`,
		},

		{
			desc: "status bad request with metadata",
			input: Status{
				Code:     http.StatusBadRequest,
				Message:  "Some message",
				Metadata: map[string]interface{}{"some": "metadata"},
			},
			expected: `{
	"code": 400,
	"message": "Some message",
	"metadata": {
		"some": "metadata"
	}
}`,
		},

		{
			desc: "status unprocessable entity with metadata",
			input: Status{
				Code:     http.StatusUnprocessableEntity,
				Message:  "Some message",
				Metadata: map[string]interface{}{"some": "metadata"},
			},
			expected: `{
	"code": 422,
	"message": "Some message",
	"metadata": {
		"some": "metadata"
	}
}`,
		},

		{
			desc: "status other with metadata",
			input: Status{
				Code:     http.StatusAccepted,
				Message:  "Accepted",
				Metadata: map[string]interface{}{"some": "metadata"},
			},
			expected: `{
	"code": 202,
	"message": "Accepted"
}`,
		},
	}

	for i, d := range data {
		result := d.input.Render()

		if string(result) != d.expected {
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
