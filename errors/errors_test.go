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
	"testing"

	go_testing "github.com/caigwatkin/go/testing"
	"github.com/pkg/errors"
)

func TestWrap(t *testing.T) {
	type input struct {
		err     error
		message string
	}
	var data = []struct {
		desc string
		input
		expected error
	}{
		{
			desc: "error",
			input: input{
				err:     errors.New("Existing message"),
				message: "A new message",
			},
			expected: errors.New("A new message: Existing message"),
		},

		{
			desc: "nil error",
			input: input{
				err:     nil,
				message: "A new message",
			},
			expected: nil,
		},

		{
			desc: "error is Status",
			input: input{
				err:     Status{},
				message: "A new message",
			},
			expected: Status{},
		},
	}

	for i, d := range data {
		result := Wrap(d.input.err, d.input.message)

		if d.expected != nil {
			if result == nil {
				t.Error(go_testing.Errorf(go_testing.Error{
					Unexpected: "result",
					Desc:       d.desc,
					At:         i,
					Expected:   d.expected,
					Result:     result,
				}))

			} else if result.Error() != d.expected.Error() {
				t.Error(go_testing.Errorf(go_testing.Error{
					Unexpected: "result.Error()",
					Desc:       d.desc,
					At:         i,
					Expected:   d.expected.Error(),
					Result:     result.Error(),
				}))
			}
			if IsStatus(result) != IsStatus(d.expected) {
				t.Error(go_testing.Errorf(go_testing.Error{
					Unexpected: "IsStatus(result)",
					Desc:       d.desc,
					At:         i,
					Expected:   IsStatus(d.expected),
					Result:     IsStatus(result),
				}))
			}

		} else if result != nil {
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

func TestWrapf(t *testing.T) {
	type input struct {
		err     error
		message string
		args    []interface{}
	}
	var data = []struct {
		desc string
		input
		expected error
	}{
		{
			desc: "error no args",
			input: input{
				err:     errors.New("Existing message"),
				message: "A new message",
				args:    nil,
			},
			expected: errors.New("A new message: Existing message"),
		},
		{
			desc: "error with args",
			input: input{
				err:     errors.New("Existing message"),
				message: "A new %s",
				args:    []interface{}{"message"},
			},
			expected: errors.New("A new message: Existing message"),
		},

		{
			desc: "nil error no args",
			input: input{
				err:     nil,
				message: "A new message",
				args:    nil,
			},
			expected: nil,
		},
		{
			desc: "nil error with args",
			input: input{
				err:     nil,
				message: "A new message",
				args:    []interface{}{"message"},
			},
			expected: nil,
		},

		{
			desc: "error is Status no args",
			input: input{
				err:     Status{},
				message: "A new message",
				args:    nil,
			},
			expected: Status{},
		},
		{
			desc: "error is Status with args",
			input: input{
				err:     Status{},
				message: "A new message",
				args:    []interface{}{"message"},
			},
			expected: Status{},
		},
	}

	for i, d := range data {
		result := Wrapf(d.input.err, d.input.message, d.input.args...)

		if d.expected != nil {
			if result == nil {
				t.Error(go_testing.Errorf(go_testing.Error{
					Unexpected: "result",
					Desc:       d.desc,
					At:         i,
					Expected:   d.expected,
					Result:     result,
				}))

			} else if result.Error() != d.expected.Error() {
				t.Error(go_testing.Errorf(go_testing.Error{
					Unexpected: "result.Error()",
					Desc:       d.desc,
					At:         i,
					Expected:   d.expected.Error(),
					Result:     result.Error(),
				}))
			}
			if IsStatus(result) != IsStatus(d.expected) {
				t.Error(go_testing.Errorf(go_testing.Error{
					Unexpected: "IsStatus(result)",
					Desc:       d.desc,
					At:         i,
					Expected:   IsStatus(d.expected),
					Result:     IsStatus(result),
				}))
			}

		} else if result != nil {
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
