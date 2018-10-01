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

package context

import (
	"context"
	"strings"
	"testing"

	go_testing "github.com/caigwatkin/go/testing"
)

func TestBackground(t *testing.T) {
	result := Background()
	expectedCorrelationID := CorrelationIDBackground
	expectedTest := false

	if CorrelationID(result) != expectedCorrelationID {
		t.Error(go_testing.Errorf(go_testing.Error{
			Unexpected: "CorrelationID(result)",
			Expected:   expectedCorrelationID,
			Result:     CorrelationID(result),
		}))
	}
	if Test(result) != expectedTest {
		t.Error(go_testing.Errorf(go_testing.Error{
			Unexpected: "Test(result)",
			Expected:   expectedTest,
			Result:     Test(result),
		}))
	}
}

func TestStartUp(t *testing.T) {
	result := StartUp()
	expectedCorrelationID := CorrelationIDStartUp
	expectedTest := false

	if CorrelationID(result) != expectedCorrelationID {
		t.Error(go_testing.Errorf(go_testing.Error{
			Unexpected: "CorrelationID(result)",
			Expected:   expectedCorrelationID,
			Result:     CorrelationID(result),
		}))
	}
	if Test(result) != expectedTest {
		t.Error(go_testing.Errorf(go_testing.Error{
			Unexpected: "Test(result)",
			Expected:   expectedTest,
			Result:     Test(result),
		}))
	}
}

func TestShutDown(t *testing.T) {
	result := ShutDown()
	expectedCorrelationID := CorrelationIDShutDown
	expectedTest := false

	if CorrelationID(result) != expectedCorrelationID {
		t.Error(go_testing.Errorf(go_testing.Error{
			Unexpected: "CorrelationID(result)",
			Expected:   expectedCorrelationID,
			Result:     CorrelationID(result),
		}))
	}
	if Test(result) != expectedTest {
		t.Error(go_testing.Errorf(go_testing.Error{
			Unexpected: "Test(result)",
			Expected:   expectedTest,
			Result:     Test(result),
		}))
	}
}

func TestNew(t *testing.T) {
	background := context.Background()
	pkgContextBackground := Background()
	pkgContextStartUp := StartUp()
	customized := context.WithValue(context.WithValue(context.Background(), keyCorrelationID, "customized"), keyTest, true)
	type expected struct {
		correlationIDSuffix string
		test                bool
	}
	var data = []struct {
		desc  string
		input context.Context
		expected
	}{
		{
			desc:  "background",
			input: background,
			expected: expected{
				correlationIDSuffix: CorrelationID(background),
				test:                Test(background),
			},
		},

		{
			desc:  "go_context background",
			input: pkgContextBackground,
			expected: expected{
				correlationIDSuffix: CorrelationID(pkgContextBackground),
				test:                Test(pkgContextBackground),
			},
		},

		{
			desc:  "go_context start up",
			input: pkgContextStartUp,
			expected: expected{
				correlationIDSuffix: CorrelationID(pkgContextStartUp),
				test:                Test(pkgContextStartUp),
			},
		},

		{
			desc:  "customized",
			input: customized,
			expected: expected{
				correlationIDSuffix: "customized",
				test:                true,
			},
		},

		{
			desc:  "nil",
			input: nil,
			expected: expected{
				correlationIDSuffix: "",
				test:                false,
			},
		},
	}

	for i, d := range data {
		result := New(d.input)

		if CorrelationID(result) == "" {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "CorrelationID(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   "NOT EMPTY STRING",
				Result:     CorrelationID(result),
			}))
		}
		if !strings.HasSuffix(CorrelationID(result), d.expected.correlationIDSuffix) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "CorrelationID(result) suffix",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.correlationIDSuffix,
				Result:     CorrelationID(result),
			}))
		}
		if Test(result) != d.expected.test {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "Test(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.test,
				Result:     Test(result),
			}))
		}
	}
}

func TestCorrelationID(t *testing.T) {
	var data = []struct {
		desc     string
		input    context.Context
		expected string
	}{
		{
			desc:     "correlationID",
			input:    context.WithValue(context.Background(), keyCorrelationID, "correlationID"),
			expected: "correlationID",
		},

		{
			desc:     "empty",
			input:    context.WithValue(context.Background(), keyCorrelationID, ""),
			expected: "",
		},

		{
			desc:     "unexpected type",
			input:    context.WithValue(context.Background(), keyCorrelationID, true),
			expected: "",
		},

		{
			desc:     "none",
			input:    context.Background(),
			expected: "",
		},
	}

	for i, d := range data {
		result := CorrelationID(d.input)

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

func TestWithCorrelationID(t *testing.T) {
	type input struct {
		ctx           context.Context
		correlationID string
	}
	var data = []struct {
		desc string
		input
		expected string
	}{
		{
			desc: "none",
			input: input{
				ctx:           context.Background(),
				correlationID: "correlationID",
			},
			expected: "correlationID",
		},

		{
			desc: "override",
			input: input{
				ctx:           context.WithValue(context.Background(), keyCorrelationID, "xxxxx"),
				correlationID: "correlationID",
			},
			expected: "correlationID",
		},
	}

	for i, d := range data {
		result := WithCorrelationID(d.input.ctx, d.input.correlationID)

		if v, ok := result.Value(keyCorrelationID).(string); !ok {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Value(keyCorrelationID).(string) ok",
				Desc:       d.desc,
				At:         i,
				Expected:   "exists",
				Result:     nil,
			}))

		} else if v != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Value(keyCorrelationID).(string)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     v,
			}))
		}
	}
}

func TestTest(t *testing.T) {
	var data = []struct {
		desc     string
		input    context.Context
		expected bool
	}{
		{
			desc:     "true",
			input:    context.WithValue(context.Background(), keyTest, true),
			expected: true,
		},

		{
			desc:     "false",
			input:    context.WithValue(context.Background(), keyTest, false),
			expected: false,
		},

		{
			desc:     "unexpected type",
			input:    context.WithValue(context.Background(), keyTest, "true"),
			expected: false,
		},

		{
			desc:     "none",
			input:    context.Background(),
			expected: false,
		},
	}

	for i, d := range data {
		result := Test(d.input)

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

func TestWithTest(t *testing.T) {
	type input struct {
		ctx  context.Context
		test bool
	}
	var data = []struct {
		desc string
		input
		expected bool
	}{
		{
			desc: "false",
			input: input{
				ctx:  context.Background(),
				test: false,
			},
			expected: false,
		},

		{
			desc: "true",
			input: input{
				ctx:  context.Background(),
				test: true,
			},
			expected: true,
		},

		{
			desc: "override",
			input: input{
				ctx:  context.WithValue(context.Background(), keyTest, true),
				test: false,
			},
			expected: false,
		},
	}

	for i, d := range data {
		result := WithTest(d.input.ctx, d.input.test)

		if v, ok := result.Value(keyTest).(bool); !ok {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Value(keyTest).(bool) ok",
				Desc:       d.desc,
				At:         i,
				Expected:   "exists",
				Result:     nil,
			}))

		} else if v != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Value(keyTest).(bool)",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected,
				Result:     v,
			}))
		}
	}
}
