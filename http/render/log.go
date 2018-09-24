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
	"net/http"

	go_log "github.com/caigwatkin/go/log"
)

func logErrorMarshallingJSONBody(ctx context.Context, logClient go_log.Client, code int, headers map[string]string) {
	logClient.Error(ctx, "Failed marshalling JSON for response body",
		go_log.FmtInt(code, "status code"),
		go_log.FmtString(http.StatusText(code), "status text"),
		go_log.FmtAny(headers, "headers"),
	)
}

func logErrorWritingBody(ctx context.Context, logClient go_log.Client, code int, headers map[string]string, body []byte) {
	logClient.Error(ctx, "Failed writing body to response",
		go_log.FmtInt(code, "status code"),
		go_log.FmtString(http.StatusText(code), "status text"),
		go_log.FmtAny(headers, "headers"),
		go_log.FmtBytes(body, "body"),
	)
}

func logInfoResponse(ctx context.Context, logClient go_log.Client, code int, headers map[string]string, lenBody int, body []byte) {
	logClient.Info(ctx, "HTTP Response",
		go_log.FmtInt(code, "status code"),
		go_log.FmtString(http.StatusText(code), "status text"),
		go_log.FmtAny(headers, "headers"),
		go_log.FmtInt(lenBody, "lenBody"),
		go_log.FmtBytes(body, "body"),
	)
}
