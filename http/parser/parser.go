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

package parser

import (
	"context"
	"io/ioutil"
	"net/http"

	go_errors "github.com/caigwatkin/go/errors"
	go_log "github.com/caigwatkin/go/log"
)

type Client interface {
	ReadRequestBody(r *http.Request) ([]byte, error)
}

func NewClient(ctx context.Context, logClient go_log.Client) Client {
	logClient.Info(ctx, "Initializing")
	logClient.Info(ctx, "Initialized")
	return client{
		logClient: logClient,
	}
}

type client struct {
	logClient go_log.Client
}

func (c client) ReadRequestBody(r *http.Request) (body []byte, err error) {
	ctx := r.Context()
	c.logClient.Info(ctx, "Reading")

	if r.Body == nil {
		c.logClient.Info(ctx, "Read", go_log.FmtBytes(body, "body"))
		return
	}

	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		err = go_errors.NewStatus(http.StatusBadRequest, "Malformed body")
		return
	}

	c.logClient.Info(ctx, "Read", go_log.FmtBytes(body, "body"))
	return
}
