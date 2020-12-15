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

package headers

import (
	"context"
	"testing"

	go_log_mock "github.com/caigwatkin/go/log/mock"
	go_testing "github.com/caigwatkin/go/testing"
)

func TestNewClient(t *testing.T) {
	var data = []struct {
		desc     string
		input    string
		expected client
	}{
		{
			desc:  "defaults",
			input: "",
			expected: client{
				correlationIdKey: correlationIdKeyDefault,
				testKey:          testKeyDefault,
			},
		},

		{
			desc:  "service",
			input: "Service-Name",
			expected: client{
				correlationIdKey: "X-Service-Name-Correlation-Id",
				testKey:          "X-Service-Name-Test",
			},
		},
	}

	for i, d := range data {
		result := NewClient(context.Background(), go_log_mock.Client, d.input)

		if result.CorrelationIdKey() != d.expected.correlationIdKey {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.CorrelationIdKey()",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.correlationIdKey,
				Result:     result.CorrelationIdKey(),
			}))
		}
		if result.TestKey() != d.expected.testKey {
			t.Error(go_testing.Errorf(go_testing.Error{
				Unexpected: "result.TestKey()",
				Desc:       d.desc,
				At:         i,
				Expected:   d.expected.testKey,
				Result:     result.TestKey(),
			}))
		}
	}
}

func TestCorrelationIdKey(t *testing.T) {
	var data = []struct {
		desc     string
		input    client
		expected string
	}{
		{
			desc: "default",
			input: client{
				correlationIdKey: correlationIdKeyDefault,
			},
			expected: correlationIdKeyDefault,
		},

		{
			desc: "foo",
			input: client{
				correlationIdKey: "foo",
			},
			expected: "foo",
		},
	}

	for i, d := range data {
		result := d.input.CorrelationIdKey()

		if result != d.expected {
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

func TestSetCorrelationIdKey(t *testing.T) {
	var data = []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "not empty service name",
			input:    "Service-Name",
			expected: "X-Service-Name-Correlation-Id",
		},

		{
			desc:     "empty service name",
			input:    "",
			expected: "X-Correlation-Id",
		},
	}

	for i, d := range data {
		c := client{}
		c.setCorrelationIdKey(d.input)
		result := c.correlationIdKey

		if result != d.expected {
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

func TestTestKey(t *testing.T) {
	var data = []struct {
		desc     string
		input    client
		expected string
	}{
		{
			desc: "default",
			input: client{
				testKey: testKeyDefault,
			},
			expected: testKeyDefault,
		},

		{
			desc: "foo",
			input: client{
				testKey: "foo",
			},
			expected: "foo",
		},
	}

	for i, d := range data {
		result := d.input.TestKey()

		if result != d.expected {
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

func TestSetTestKey(t *testing.T) {
	var data = []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "not empty service name",
			input:    "Service-Name",
			expected: "X-Service-Name-Test",
		},

		{
			desc:     "empty service name",
			input:    "",
			expected: "X-Test",
		},
	}

	for i, d := range data {
		c := client{}
		c.setTestKey(d.input)
		result := c.testKey

		if result != d.expected {
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
