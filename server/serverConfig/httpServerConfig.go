/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package serverConfig

import (
	"net/http"

	"github.com/rbell/toolchest/server/httpMiddleware"
)

type HttpServerConfig struct {
	Port       string
	routes     map[string]map[string]http.HandlerFunc
	middleware httpMiddleware.HttpHandlerMiddleware
}

func (c *HttpServerConfig) SetMiddleware(middleware httpMiddleware.HttpHandlerMiddleware) {
	c.middleware = middleware
}

func (c *HttpServerConfig) GetMiddleware() httpMiddleware.HttpHandlerMiddleware {
	return c.middleware
}

func (c *HttpServerConfig) AddRoute(method, path string, handler http.HandlerFunc) {
	if c.routes == nil {
		c.routes = map[string]map[string]http.HandlerFunc{}
	}
	if c.routes[method] == nil {
		c.routes[method] = map[string]http.HandlerFunc{}
	}
	c.routes[method][path] = handler
}

func (c *HttpServerConfig) GetRoutes() map[string]map[string]http.HandlerFunc {
	return c.routes
}
