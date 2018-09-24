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

package middleware

import (
	"net/http"
	"strings"
	"time"

	go_context "github.com/caigwatkin/go/context"
	go_headers "github.com/caigwatkin/go/http/headers"
	go_log "github.com/caigwatkin/go/log"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
)

// Defaults adds middleware defaults to the router
func Defaults(router *chi.Mux, headersClient go_headers.Client, logClient go_log.Client, excludePathsForLogInfoRequests []string) {
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(time.Second * 30))
	router.Use(middleware.URLFormat)
	router.Use(populateContext(headersClient))
	router.Use(logInfoRequests(logClient, excludePathsForLogInfoRequests))
	router.Use(middleware.DefaultCompress)
}

func populateContext(headersClient go_headers.Client) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := go_context.WithCorrelationID(r.Context(), uuid.New().String())
			if v, ok := r.Header[headersClient.CorrelationIDKey()]; ok {
				ctx = go_context.WithCorrelationIDAppend(ctx, strings.Join(v, ","))
			}
			var t bool
			if v, ok := r.Header[headersClient.TestKey()]; ok {
				t = strings.Join(v, ",") == go_headers.TestValDefault
			}
			ctx = go_context.WithTest(ctx, t)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func logInfoRequests(logClient go_log.Client, excludePaths []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			url := r.URL.String()
			var exclude bool
			for _, v := range excludePaths {
				if url == v {
					exclude = true
					break
				}
			}
			if !exclude {
				logClient.Info(r.Context(), "HTTP Request",
					go_log.FmtString(r.URL.String(), "URL"),
					go_log.FmtString(r.Method, "Method"),
					go_log.FmtAny(r.Header, "Header"),
				)
			}
			next.ServeHTTP(w, r)
		})
	}
}
