/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package httpMiddleware

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/rbell/toolchest/server/internal/sharedTypes"
)

type RequestLogDetail struct {
	Method string
	URL    string
	Body   string
}

func LogRequest(logger sharedTypes.LogPublisher) HttpHandlerMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logRequest(logger, r)
			next(w, r)
		}
	}
}

func logRequest(logger sharedTypes.LogPublisher, r *http.Request) {
	body, _ := getRequestBodyAsString(r)
	logger.InfoContext(r.Context(), "Request received", slog.Any("request", RequestLogDetail{Method: r.Method, URL: r.URL.Path, Body: body}))
}

func getRequestBodyAsString(r *http.Request) (string, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	//nolint:errcheck // skip error in defer
	defer r.Body.Close()

	bodyString := string(bodyBytes)
	return bodyString, nil
}
