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
	return newStatus(1, code, message, nil)
}

// Statusf with code and formatted message
func Statusf(code int, format string, args ...interface{}) Status {
	return newStatus(1, code, fmt.Sprintf(format, args...), nil)
}

// NewStatusWithItems with code, message, and items
//
// Items can be useful to add extra context to the error
func NewStatusWithItems(code int, message string, items []Item) Status {
	return newStatus(1, code, message, items)
}

func newStatus(atSkip, code int, message string, items []Item) Status {
	pc, _, lineNumber, _ := runtime.Caller(atSkip + 1)
	s := Status{
		At:      fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineNumber),
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
	return fmt.Sprintf("%q: %q, %q: %q, %q: %q", "code", s.Code, "message", s.Message, "at", s.At)
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
