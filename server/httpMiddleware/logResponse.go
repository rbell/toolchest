/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package httpMiddleware

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rbell/toolchest/server/internal/sharedTypes"
)

type ResponseLogDetail struct {
	Method     string
	URL        string
	StatusCode string
	Response   string
}

func LogResponse(logger sharedTypes.LogPublisher) HttpHandlerMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// wrap the response writer to capture the response for logging
			rww := newResponseWriterWrapper(w)
			defer logResponse(logger, rww, r)

			// call the next middleware in the chain with the wrapped response writer
			next(rww, r)
		}
	}
}

func logResponse(logger sharedTypes.LogPublisher, rww ResponseWriterWrapper, r *http.Request) {
	logger.InfoContext(r.Context(), "Response Sent", slog.Any("response", ResponseLogDetail{
		Method:     r.Method,
		URL:        r.URL.Path,
		StatusCode: rww.getStatusCode(),
		Response:   rww.getResponseBody(),
	}))
}

// ResponseWriterWrapper struct is used to log the response
type ResponseWriterWrapper struct {
	w          *http.ResponseWriter
	body       *bytes.Buffer
	statusCode *int
}

// NewResponseWriterWrapper static function creates a wrapper for the http.ResponseWriter
func newResponseWriterWrapper(w http.ResponseWriter) ResponseWriterWrapper {
	var buf bytes.Buffer
	var statusCode int = 200
	return ResponseWriterWrapper{
		w:          &w,
		body:       &buf,
		statusCode: &statusCode,
	}
}

func (rww ResponseWriterWrapper) Write(buf []byte) (int, error) {
	rww.body.Write(buf)
	return (*rww.w).Write(buf)
}

// Header function overwrites the http.ResponseWriter Header() function
func (rww ResponseWriterWrapper) Header() http.Header {
	return (*rww.w).Header()
}

// WriteHeader function overwrites the http.ResponseWriter WriteHeader() function
func (rww ResponseWriterWrapper) WriteHeader(statusCode int) {
	(*rww.statusCode) = statusCode
	(*rww.w).WriteHeader(statusCode)
}

func (rww ResponseWriterWrapper) getResponseBody() string {
	var buf bytes.Buffer
	buf.WriteString(rww.body.String())
	return buf.String()
}

func (rww ResponseWriterWrapper) getStatusCode() string {
	return fmt.Sprintf("%d", *(rww.statusCode))
}
