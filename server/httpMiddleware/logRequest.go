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
)

type RequestLogDetail struct {
	Method string
	URL    string
	Body   string
}

func LogRequest(logger *slog.Logger) HttpHandlerMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			body, _ := getRequestBodyAsString(r)
			ctx := r.Context()
			if ctx == nil {
				logger.Info("Request received", slog.Any("request", RequestLogDetail{Method: r.Method, URL: r.URL.Path, Body: body}))
			} else {
				logger.InfoContext(ctx, "Request received", slog.Any("request", RequestLogDetail{Method: r.Method, URL: r.URL.Path, Body: body}))
			}
			next(w, r)
		}
	}
}

func getRequestBodyAsString(r *http.Request) (string, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	// It's a good practice to close the body to avoid resource leaks
	defer r.Body.Close()

	bodyString := string(bodyBytes)
	return bodyString, nil
}
