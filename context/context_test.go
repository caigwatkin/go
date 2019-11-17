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
	expectedCorrelationId := CorrelationIdBackground
	expectedTest := false

	if CorrelationId(result) != expectedCorrelationId {
		t.Error(go_testing.Errorf(go_testing.Error{
			Unexpected: "CorrelationId(result)",
			Expected:   expectedCorrelationId,
			Result:     CorrelationId(result),
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
	expectedCorrelationId := CorrelationIdStartUp
	expectedTest := false

	if CorrelationId(result) != expectedCorrelationId {
		t.Error(go_testing.Errorf(go_testing.Error{
			Unexpected: "CorrelationId(result)",
			Expected:   expectedCorrelationId,
			Result:     CorrelationId(result),
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
	expectedCorrelationId := CorrelationIdShutDown
	expectedTest := false

	if CorrelationId(result) != expectedCorrelationId {
		t.Error(go_testing.Errorf(go_testing.Error{
			Unexpected: "CorrelationId(result)",
			Expected:   expectedCorrelationId,
			Result:     CorrelationId(result),
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
	customized := context.WithValue(context.WithValue(context.Background(), keyCorrelationId, "customized"), keyTest, true)
	type expected struct {
		correlationIdSuffix string
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
				correlationIdSuffix: CorrelationId(background),
				test:                Test(background),
			},
		},

		{
			desc:  "go_context background",
			input: pkgContextBackground,
			expected: expected{
				correlationIdSuffix: CorrelationId(pkgContextBackground),
				test:                Test(pkgContextBackground),
			},
		},

		{
			desc:  "go_context start up",
			input: pkgContextStartUp,
			expected: expected{
				correlationIdSuffix: CorrelationId(pkgContextStartUp),
				test:                Test(pkgContextStartUp),
			},
		},

		{
			desc:  "customized",
			input: customized,
			expected: expected{
				correlationIdSuffix: "customized",
				test:                true,
			},
		},

		{
			desc:  "nil",
			input: nil,
			expected: expected{
				correlationIdSuffix: "",
				test:                false,
			},
		},
	}

	for i, d := range data {
		result := New(d.input)

		if CorrelationId(result) == "" {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "CorrelationId(result)",
				Desc:       d.desc,
				At:         i,
				Expected:   "NOT EMPTY STRING",
				Result:     CorrelationId(result),
			}))
		}
		if !strings.HasSuffix(CorrelationId(result), d.expected.correlationIdSuffix) {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "CorrelationId(result) suffix",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.correlationIdSuffix,
				Result:     CorrelationId(result),
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

func TestCorrelationId(t *testing.T) {
	var data = []struct {
		desc     string
		input    context.Context
		expected string
	}{
		{
			desc:     "correlationId",
			input:    context.WithValue(context.Background(), keyCorrelationId, "correlationId"),
			expected: "correlationId",
		},

		{
			desc:     "empty",
			input:    context.WithValue(context.Background(), keyCorrelationId, ""),
			expected: "",
		},

		{
			desc:     "unexpected type",
			input:    context.WithValue(context.Background(), keyCorrelationId, true),
			expected: "",
		},

		{
			desc:     "none",
			input:    context.Background(),
			expected: "",
		},
	}

	for i, d := range data {
		result := CorrelationId(d.input)

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

func TestWithCorrelationId(t *testing.T) {
	type input struct {
		ctx           context.Context
		correlationId string
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
				correlationId: "correlationId",
			},
			expected: "correlationId",
		},

		{
			desc: "override",
			input: input{
				ctx:           context.WithValue(context.Background(), keyCorrelationId, "xxxxx"),
				correlationId: "correlationId",
			},
			expected: "correlationId",
		},
	}

	for i, d := range data {
		result := WithCorrelationId(d.input.ctx, d.input.correlationId)

		if v, ok := result.Value(keyCorrelationId).(string); !ok {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Value(keyCorrelationId).(string) ok",
				Desc:       d.desc,
				At:         i,
				Expected:   "exists",
				Result:     nil,
			}))

		} else if v != d.expected {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.Value(keyCorrelationId).(string)",
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
