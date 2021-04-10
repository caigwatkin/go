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
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

// Status data model
//
// Implements error interface
type Status struct {
	At      string
	Cause   error
	Code    int
	Message string
	Items   []Item
}

type Item struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

// NewStatus with code and message
func NewStatus(code int, message string) Status {
	return newStatus(1, nil, code, message, nil)
}

// Statusf with code and formatted message
func Statusf(code int, format string, args ...interface{}) Status {
	return newStatus(1, nil, code, fmt.Sprintf(format, args...), nil)
}

// NewStatus with cause, code, and message
//
// Cause can be useful to record an error that caused a Status to be created
func NewStatusWithCause(cause error, code int, message string) Status {
	return newStatus(1, cause, code, message, nil)
}

// NewStatusWithItems with code, message, and items
//
// Items can be useful to add extra context to the error
func NewStatusWithItems(code int, message string, items []Item) Status {
	return newStatus(1, nil, code, message, items)
}

func newStatus(atSkip int, cause error, code int, message string, items []Item) Status {
	pc, _, lineNumber, _ := runtime.Caller(atSkip + 1)
	s := Status{
		At:      fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineNumber),
		Cause:   cause,
		Code:    code,
		Message: http.StatusText(code),
		Items:   items,
	}
	if message != "" {
		s.Message = fmt.Sprintf("%s: %s", s.Message, message)
	}
	return s
}

// StatusCode returns the Status code if err is a Status, zero if err is not a Status
func StatusCode(err error) int {
	if s, ok := err.(Status); ok {
		return s.Code
	}
	return 0
}

// IsStatus returns true if err is a Status
func IsStatus(err error) bool {
	_, ok := err.(Status)
	return ok
}

// Error so that Status objects can be treated as errors
func (s Status) Error() string {
	e := fmt.Sprintf("Code: %d, Message: %q, At: %q, Items: %v", s.Code, s.Message, s.At, s.Items)
	if s.Cause != nil {
		e = fmt.Sprintf("%s, Cause: %+v", e, s.Cause)
	}
	return e
}

// Render items
func (s Status) RenderItems() []byte {
	if len(s.Items) == 0 {
		return nil
	}
	b, err := json.Marshal(s.Items)
	if err != nil {
		return nil
	}
	return b
}
