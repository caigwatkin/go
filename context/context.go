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
	"fmt"

	"github.com/google/uuid"
)

// Correlation ID enums
const (
	CorrelationIdBackground = "BACKGROUND"
	CorrelationIdStartUp    = "START_UP"
	CorrelationIdShutDown   = "SHUT_DOWN"
)

// Background returns a new background context with background correlation ID enum
func Background() context.Context {
	return WithCorrelationId(context.Background(), CorrelationIdBackground)
}

// StartUp returns a new background context with start up correlation ID enum
func StartUp() context.Context {
	return WithCorrelationId(context.Background(), CorrelationIdStartUp)
}

// ShutDown returns a new background context with shut down correlation ID enum
func ShutDown() context.Context {
	return WithCorrelationId(context.Background(), CorrelationIdShutDown)
}

// New context with correlation ID of ctx with newly appended ctx, test value of ctx, and other defaults
func New(ctx context.Context) context.Context {
	c := WithCorrelationId(context.Background(), uuid.New().String())
	if ctx != nil {
		c = WithCorrelationIdAppend(c, CorrelationId(ctx))
		c = WithTest(c, Test(ctx))
	}
	return c
}

type key int

const (
	keyCorrelationId key = iota
	keyTest          key = iota
)

// CorrelationId returns correlation ID value of ctx
func CorrelationId(ctx context.Context) string {
	if v, ok := ctx.Value(keyCorrelationId).(string); ok {
		return v
	}
	return ""
}

// WithCorrelationId returns a new context with correlation ID value
func WithCorrelationId(ctx context.Context, correlationId string) context.Context {
	return context.WithValue(ctx, keyCorrelationId, correlationId)
}

// WithCorrelationIdAppend returns a new context with correlation ID value appended to existing correlation ID value if one exists
func WithCorrelationIdAppend(ctx context.Context, correlationId string) context.Context {
	if cId := CorrelationId(ctx); cId != CorrelationIdBackground {
		correlationId = fmt.Sprintf("%s,%s", cId, correlationId)
	}
	return WithCorrelationId(ctx, correlationId)
}

// Test returns test value of ctx
func Test(ctx context.Context) bool {
	if v, ok := ctx.Value(keyTest).(bool); ok {
		return v
	}
	return false
}

// WithTest returns a new context with test value
func WithTest(ctx context.Context, test bool) context.Context {
	return context.WithValue(ctx, keyTest, test)
}
