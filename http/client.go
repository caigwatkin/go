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

package http

import (
	"context"
	"net/http"

	go_headers "github.com/caigwatkin/go/http/headers"
	go_middleware "github.com/caigwatkin/go/http/middleware"
	go_render "github.com/caigwatkin/go/http/render"
	go_log "github.com/caigwatkin/go/log"
	"github.com/go-chi/chi"
)

type Client interface {
	RenderContentJSON(ctx context.Context, w http.ResponseWriter, body []byte)
	RenderHealth(ctx context.Context, w http.ResponseWriter, serviceName string)
	RenderErrorOrStatus(ctx context.Context, w http.ResponseWriter, err error)
	MiddlewareDefaults(r *chi.Mux, excludePathsForLogInfoRequests []string)
}

type client struct {
	headersClient go_headers.Client
	logClient     go_log.Client
}

// NewClient for http
//
// Service name should be in canonical case as it is used for custom response headers
// Use an empty string to use default keys
func NewClient(logClient go_log.Client, serviceNameForHeaders string) Client {
	return client{
		headersClient: go_headers.NewClient(serviceNameForHeaders),
		logClient:     logClient,
	}
}

// RenderContentJSON in response
func (c client) RenderContentJSON(ctx context.Context, w http.ResponseWriter, body []byte) {
	go_render.ContentJSON(ctx, c.headersClient, c.logClient, w, body)
}

// RenderHealth in response
func (c client) RenderHealth(ctx context.Context, w http.ResponseWriter, serviceName string) {
	go_render.Health(ctx, c.headersClient, c.logClient, w, serviceName)
}

// RenderErrorOrStatus in response
func (c client) RenderErrorOrStatus(ctx context.Context, w http.ResponseWriter, err error) {
	go_render.ErrorOrStatus(ctx, c.headersClient, c.logClient, w, err)
}

// MiddlewareDefaults for request handling
func (c client) MiddlewareDefaults(r *chi.Mux, excludePathsForLogInfoRequests []string) {
	go_middleware.Defaults(r, c.headersClient, c.logClient, excludePathsForLogInfoRequests)
}
