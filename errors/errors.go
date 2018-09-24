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
	"github.com/pkg/errors"
)

// New error with stack trace
func New(message string) error {
	return errors.New(message)
}

// Errorf with stack trace and formatted message
func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

// Wrap error with new message
//
// If err is a Status, it will not be wrapped
func Wrap(err error, message string) error {
	if IsStatus(err) {
		return err
	}
	return errors.Wrap(err, message)
}

// Wrapf error with new formatted message
//
// If err is a Status, it will not be wrapped
func Wrapf(err error, format string, args ...interface{}) error {
	if IsStatus(err) {
		return err
	}
	return errors.Wrapf(err, format, args...)
}
