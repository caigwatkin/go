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

	go_context "github.com/caigwatkin/go/context"
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

// Created writes location header to response writer with status code Created
func Created(ctx context.Context, headersClient go_headers.Client, logClient go_log.Client, w http.ResponseWriter, location string) {
	headers := setHeadersInclDefaults(ctx, headersClient, w, map[string]string{
		"Location": location,
	})
	code := http.StatusCreated
	w.WriteHeader(code)
	logInfoResponse(ctx, logClient, code, headers, 0, nil)
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
	logStatus(ctx, logClient, s)
	body := s.Render()
	body = append(body, byte('\n'))
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
	logError(ctx, logClient, err)
	h := setHeadersInclDefaults(ctx, headersClient, w, nil)
	code := http.StatusInternalServerError
	w.WriteHeader(code)
	logInfoResponse(ctx, logClient, code, h, 0, nil)
}

func setHeadersInclDefaults(ctx context.Context, headersClient go_headers.Client, w http.ResponseWriter, headers map[string]string) map[string]string {
	h := map[string]string{
		headersClient.CorrelationIDKey(): go_context.CorrelationID(ctx),
	}
	if go_context.Test(ctx) {
		h[headersClient.TestKey()] = go_headers.TestValDefault
	}
	for k, v := range headers {
		h[k] = v
	}
	for k, v := range h {
		w.Header().Set(k, v)
	}
	return h
}
