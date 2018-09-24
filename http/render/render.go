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

package render

import (
	"context"
	"encoding/json"
	"net/http"

	go_errors "github.com/caigwatkin/go/errors"
	go_headers "github.com/caigwatkin/go/http/headers"
	go_log "github.com/caigwatkin/go/log"
)

// ContentJSON writes JSON bytes to the response writer with status code OK
func ContentJSON(ctx context.Context, headersClient go_headers.Client, logClient go_log.Client, w http.ResponseWriter, body []byte) {
	headers := setHeadersInclDefaults(ctx, headersClient, w, map[string]string{
		"Content-Type": "application/json",
	})
	body = append(body, byte('\n'))
	code := http.StatusOK
	w.WriteHeader(code)
	lenBody, err := w.Write(body)
	if err != nil {
		logErrorWritingBody(ctx, logClient, code, headers, body)
		return
	}
	logInfoResponse(ctx, logClient, code, headers, lenBody, body)
}

// ContentJSON writes a small JSON body to the response writer with status code OK
//
// Used for health check endpoints to ensure API is serving
func Health(ctx context.Context, headersClient go_headers.Client, logClient go_log.Client, w http.ResponseWriter, serviceName string) {
	headers := setHeadersInclDefaults(ctx, headersClient, w, map[string]string{
		"Content-Type": "application/json",
	})
	body, err := json.MarshalIndent(map[string]string{
		"service": serviceName,
		"status":  "OK",
	}, "", "\t")
	if err != nil {
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		logErrorMarshallingJSONBody(ctx, logClient, code, headers)
		return
	}
	body = append(body, byte('\n'))
	code := http.StatusOK
	w.WriteHeader(code)
	if _, err := w.Write(body); err != nil {
		logErrorWritingBody(ctx, logClient, code, headers, body)
		return
	}
}

// Status writes a go_errors.Status as JSON to the response writer with status code Status.Code
func Status(ctx context.Context, headersClient go_headers.Client, logClient go_log.Client, w http.ResponseWriter, s go_errors.Status) {
	h := setHeadersInclDefaults(ctx, headersClient, w, map[string]string{
		"Content-Type": "application/json",
	})
	body := s.Render()
	w.WriteHeader(s.Code)
	lenBody, err := w.Write(body)
	if err != nil {
		logErrorWritingBody(ctx, logClient, s.Code, h, body)
		return
	}
	logInfoResponse(ctx, logClient, s.Code, h, lenBody, body)
}

// ErrorOrStatus wraps Status if error is a Status, otherwise writes status code Internal Server Error
func ErrorOrStatus(ctx context.Context, headersClient go_headers.Client, logClient go_log.Client, w http.ResponseWriter, err error) {
	if v, ok := err.(go_errors.Status); ok {
		Status(ctx, headersClient, logClient, w, v)
		return
	}
	h := setHeadersInclDefaults(ctx, headersClient, w, nil)
	code := http.StatusInternalServerError
	w.WriteHeader(code)
	logInfoResponse(ctx, logClient, code, h, 0, nil)
}
