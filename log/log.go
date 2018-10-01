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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	go_context "github.com/caigwatkin/go/context"
)

type Client interface {
	Debug(ctx context.Context, message string, fields ...Field)
	Info(ctx context.Context, message string, fields ...Field)
	Notice(ctx context.Context, message string, fields ...Field)
	Warn(ctx context.Context, message string, fields ...Field)
	Error(ctx context.Context, message string, fields ...Field)
	Fatal(ctx context.Context, message string, fields ...Field)
}

type client struct {
	debug       bool
	loggerDebug *log.Logger
	loggerInfo  *log.Logger
	loggerWarn  *log.Logger
	loggerError *log.Logger
	loggerFatal *log.Logger
}

// NewClient for logging
func NewClient(enableDebug bool) Client {
	return client{
		debug:       enableDebug,
		loggerDebug: log.New(os.Stdout, fmt.Sprintf("\x1b[%dmDEBUG ", green), log.Ldate|log.Ltime|log.Lmicroseconds),
		loggerInfo:  log.New(os.Stdout, fmt.Sprintf("\x1b[%dmINFO  ", cyan), log.Ldate|log.Ltime|log.Lmicroseconds),
		loggerWarn:  log.New(os.Stderr, fmt.Sprintf("\x1b[%dmWARN  ", yellow), log.Ldate|log.Ltime|log.Lmicroseconds),
		loggerError: log.New(os.Stderr, fmt.Sprintf("\x1b[%dmERROR ", red), log.Ldate|log.Ltime|log.Lmicroseconds),
		loggerFatal: log.New(os.Stderr, fmt.Sprintf("\x1b[%dmFATAL ", red), log.Ldate|log.Ltime|log.Lmicroseconds),
	}
}

// Debug log at debug level
func (c client) Debug(ctx context.Context, message string, fields ...Field) {
	if c.debug {
		c.output(ctx, severityDebug, message, fields)
	}
}

// Info log at info level
func (c client) Info(ctx context.Context, message string, fields ...Field) {
	c.output(ctx, severityInfo, message, fields)
}

// Notice log at notice level
func (c client) Notice(ctx context.Context, message string, fields ...Field) {
	c.output(ctx, severityNotice, message, fields)
}

// Warn log at warn level
func (c client) Warn(ctx context.Context, message string, fields ...Field) {
	c.output(ctx, severityWarn, message, fields)
}

// Error log at error level
func (c client) Error(ctx context.Context, message string, fields ...Field) {
	c.output(ctx, severityError, message, fields)
}

// Fatal log at fatal level
func (c client) Fatal(ctx context.Context, message string, fields ...Field) {
	c.output(ctx, severityFatal, message, fields)
}

// Field to log
type Field string

// FmtAny using JSON marshaller with indenting
//
// Use this to format for logging any object that can be JSON unmarshalled
// Type and value will be logged
// If JSON unmarshalling fails, the value will not be logged
func FmtAny(value interface{}, name string) Field {
	if value == nil {
		return Field(fmt.Sprintf("%q: null", name))
	}
	blob, err := json.MarshalIndent(value, "\t\t", "\t")
	if err != nil {
		return Field(fmt.Sprintf("%q: {\n\t\t\"type\": %q,\n\t\t\"value\": \"NOT JSON MARSHALLABLE\"\n\t}", name, reflect.TypeOf(value)))
	}
	return Field(fmt.Sprintf("%q: {\n\t\t\"type\": %q,\n\t\t\"value\": %s\n\t}", name, reflect.TypeOf(value), blob))
}

// FmtAnys using JSON marshaller with indenting
//
// Use this to format for logging any object that can be JSON unmarshalled
// Type and value will be logged
// If JSON unmarshalling fails, the value will not be logged
func FmtAnys(values []interface{}, name string) Field {
	var vals []interface{}
	for _, v := range values {
		if v == nil {
			vals = append(vals, "null")
			continue
		}
		var s string
		if blob, err := json.MarshalIndent(v, "\t\t\t", "\t"); err != nil {
			s = fmt.Sprintf("{\n\t\t\t\"type\": %q,\n\t\t\t\"value\": \"NOT JSON MARSHALLABLE\"\n\t\t}", reflect.TypeOf(v))
		} else {
			s = fmt.Sprintf("{\n\t\t\t\"type\": %q,\n\t\t\t\"value\": %s\n\t\t}", reflect.TypeOf(v), blob)
		}
		vals = append(vals, s)
	}
	return fmtSlice(vals, name, "%s")
}

// FmtBool as name/value pair for logging
func FmtBool(value bool, name string) Field {
	return Field(fmt.Sprintf("%q: %t", name, value))
}

// FmtBools as name/[values] pair for logging
func FmtBools(values []bool, name string) Field {
	var vals []interface{}
	for _, v := range values {
		vals = append(vals, v)
	}
	return fmtSlice(vals, name, "%t")
}

// FmtByte as name/value pair for logging
func FmtByte(value byte, name string) Field {
	return Field(fmt.Sprintf("%q: %q", name, value))
}

// FmtBytes as name/value pair for logging
func FmtBytes(value []byte, name string) Field {
	return Field(fmt.Sprintf("%q: %q", name, value))
}

// FmtDuration as name/value pair for logging
func FmtDuration(value time.Duration, name string) Field {
	return Field(fmt.Sprintf("%q: %q", name, value))
}

// FmtDurations as name/[values] pair for logging
func FmtDurations(values []time.Duration, name string) Field {
	var vals []interface{}
	for _, v := range values {
		vals = append(vals, v)
	}
	return fmtSlice(vals, name, "%q")
}

// FmtError as name/value pair for logging
func FmtError(err error) Field {
	if err == nil {
		return Field("\"error\": null")
	}
	friendly := fmt.Sprintf("%s", err)
	trace := fmt.Sprintf("%+v", err)
	if trace == friendly {
		return Field(fmt.Sprintf("\"error\": %q", friendly))
	}
	return Field(fmt.Sprintf("\"error\": {\n\t\t\"friendly\": %q,\n\t\t\"trace\": %s\n\t}", friendly, trace))
}

// FmtFloat32 as name/value pair for logging
func FmtFloat32(value float32, name string) Field {
	return Field(fmt.Sprintf("%q: %.5f", name, value))
}

// FmtFloat32s as name/[values] pair for logging
func FmtFloat32s(values []float32, name string) Field {
	var vals []interface{}
	for _, v := range values {
		vals = append(vals, v)
	}
	return fmtSlice(vals, name, "%.5f")
}

// FmtFloat64 as name/value pair for logging
func FmtFloat64(value float64, name string) Field {
	return Field(fmt.Sprintf("%q: %.10f", name, value))
}

// FmtFloat64s as name/[values] pair for logging
func FmtFloat64s(values []float64, name string) Field {
	var vals []interface{}
	for _, v := range values {
		vals = append(vals, v)
	}
	return fmtSlice(vals, name, "%.10f")
}

// FmtInt as name/value pair for logging
func FmtInt(value int, name string) Field {
	return Field(fmt.Sprintf("%q: %d", name, value))
}

// FmtInts as name/[values] pair for logging
func FmtInts(values []int, name string) Field {
	var vals []interface{}
	for _, v := range values {
		vals = append(vals, v)
	}
	return fmtSlice(vals, name, "%d")
}

// FmtInt32 as name/value pair for logging
func FmtInt32(value int32, name string) Field {
	return Field(fmt.Sprintf("%q: %d", name, value))
}

// FmtInt32s as name/[values] pair for logging
func FmtInt32s(values []int32, name string) Field {
	var vals []interface{}
	for _, v := range values {
		vals = append(vals, v)
	}
	return fmtSlice(vals, name, "%d")
}

// FmtInt64 as name/value pair for logging
func FmtInt64(value int64, name string) Field {
	return Field(fmt.Sprintf("%q: %d", name, value))
}

// FmtInt64s as name/[values] pair for logging
func FmtInt64s(values []int64, name string) Field {
	var vals []interface{}
	for _, v := range values {
		vals = append(vals, v)
	}
	return fmtSlice(vals, name, "%d")
}

// FmtString as name/value pair for logging
func FmtString(value string, name string) Field {
	return Field(fmt.Sprintf("%q: %q", name, value))
}

// FmtStrings as name/[values] pair for logging
func FmtStrings(values []string, name string) Field {
	var vals []interface{}
	for _, v := range values {
		vals = append(vals, v)
	}
	return fmtSlice(vals, name, "%q")
}

// FmtTime as name/value pair for logging
func FmtTime(value time.Time, name string) Field {
	return Field(fmt.Sprintf("%q: %q", name, value.Format(time.RFC3339Nano)))
}

// FmtTimes as name/[values] pair for logging
func FmtTimes(values []time.Time, name string) Field {
	var vals []interface{}
	for _, v := range values {
		vals = append(vals, v.Format(time.RFC3339Nano))
	}
	return fmtSlice(vals, name, "%q")
}

func fmtSlice(values []interface{}, name, format string) Field {
	if len(values) == 0 {
		return Field(fmt.Sprintf("%q: []", name))
	}
	f := make([]string, len(values))
	for i, v := range values {
		f[i] = fmt.Sprintf(format, v)
	}
	return Field(fmt.Sprintf("%q: [\n\t\t%s\n\t]", name, strings.Join(f, ",\n\t\t")))
}

const (
	red    = 31
	green  = 32
	yellow = 33
	cyan   = 36

	severityDebug  = iota
	severityInfo   = iota
	severityNotice = iota // Flush, if applicable
	severityWarn   = iota
	severityError  = iota
	severityFatal  = iota // Flush, if applicable
)

func (c client) output(ctx context.Context, severity int, message string, fields []Field) {
	line, funcName := runtimeLineAndFuncName(2)
	message = fmtLog(go_context.CorrelationID(ctx), message, funcName, line, fields)
	switch severity {
	case severityDebug:
		c.loggerDebug.Println(message)
	case severityInfo, severityNotice:
		c.loggerInfo.Println(message)
	case severityWarn:
		c.loggerWarn.Println(message)
	case severityError:
		c.loggerError.Println(message)
	case severityFatal:
		c.loggerFatal.Panicln(message)
	}
}

func runtimeLineAndFuncName(skip int) (int, string) {
	pc, _, line, _ := runtime.Caller(skip + 1)
	funcName := runtime.FuncForPC(pc).Name()
	return line, funcName
}

func fmtLog(message, correlationID, funcName string, line int, fields []Field) string {
	return fmt.Sprintf("%s %s %s:%d %s\x1b[0m", message, correlationID, funcName, line, fmtFields(fields))
}

func fmtFields(fields []Field) string {
	if len(fields) == 0 {
		return ""
	}
	var fs []string
	for _, field := range fields {
		fs = append(fs, string(field))
	}
	return fmt.Sprintf("{\n\t%s\n}", strings.Join(fs, ",\n\t"))
}
