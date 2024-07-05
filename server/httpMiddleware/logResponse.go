/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package httpMiddleware

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type ResponseLogDetail struct {
	Method     string
	URL        string
	StatusCode string
	Response   string
}

func LogResponse(logger *slog.Logger) HttpHandlerMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// wrap the response writer to capture the response for logging
			rww := newResponseWriterWrapper(w)

			defer func() {
				// log the response
				ctx := r.Context()
				if ctx == nil {
					logger.Info("Response Sent", slog.Any("response", ResponseLogDetail{
						Method:     r.Method,
						URL:        r.URL.Path,
						StatusCode: rww.getStatusCode(),
						Response:   rww.getResponseBody(),
					}))
				} else {
					logger.InfoContext(ctx, "Response Sent", slog.Any("response", ResponseLogDetail{
						Method:     r.Method,
						URL:        r.URL.Path,
						StatusCode: rww.getStatusCode(),
						Response:   rww.getResponseBody(),
					}))
				}
			}()

			// call the next middleware in the chain with the wrapped response writer
			next(rww, r)
		}
	}
}

func getResponseAsString(r *http.Request) (string, error) {
	bodyBytes, err := io.ReadAll(r.Response.Body)
	if err != nil {
		return "", err
	}
	// It's a good practice to close the body to avoid resource leaks
	defer r.Body.Close()

	bodyString := string(bodyBytes)
	return bodyString, nil
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

func (rww ResponseWriterWrapper) getHeaders() string {
	var buf bytes.Buffer

	for k, v := range (*rww.w).Header() {
		buf.WriteString(fmt.Sprintf("%s: %v", k, v))
	}

	return buf.String()
}

func (rww ResponseWriterWrapper) getResponseBody() string {
	var buf bytes.Buffer
	buf.WriteString(rww.body.String())
	return buf.String()
}

func (rww ResponseWriterWrapper) getStatusCode() string {
	return fmt.Sprintf("%d", *(rww.statusCode))
}
