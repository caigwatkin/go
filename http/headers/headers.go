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
	"fmt"
)

// Client interface
type Client interface {
	CorrelationIdKey() string
	TestKey() string
}

type client struct {
	correlationIdKey string
	testKey          string
}

const (
	correlationIdKeyDefault = "X-Correlation-Id"
	correlationIdKeyFormat  = "X-%s-Correlation-Id"

	testKeyDefault = "X-Test"
	testKeyFormat  = "X-%s-Test"
	TestValDefault = "5c1bca85-9e09-4af4-96ac-7f353265838c" // This can be stored and used unencrypted, as anything running in test mode should be safe enough that it doesn't matter who knows it
)

// NewClient with defaults
//
// Service name should be in canonical case
// Use an empty string to use default keys
func NewClient(serviceName string) Client {
	var c client
	c.setCorrelationIdKey(serviceName)
	c.setTestKey(serviceName)
	return &c
}

// CorrelationIdKey returns the correlation ID header key
func (c client) CorrelationIdKey() string {
	return c.correlationIdKey
}

func (c *client) setCorrelationIdKey(serviceName string) {
	if serviceName == "" {
		c.correlationIdKey = correlationIdKeyDefault
		return
	}
	c.correlationIdKey = fmt.Sprintf(correlationIdKeyFormat, serviceName)
}

// TestKey returns the test header key
func (c client) TestKey() string {
	return c.testKey
}

func (c *client) setTestKey(serviceName string) {
	if serviceName == "" {
		c.testKey = testKeyDefault
		return
	}
	c.testKey = fmt.Sprintf(testKeyFormat, serviceName)
}
