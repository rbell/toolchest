/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package serverConfig

import "github.com/julienschmidt/httprouter"

type HttpServerConfig struct {
	Port   string
	routes map[string]map[string]httprouter.Handle
}

func (c *HttpServerConfig) AddRoute(method, path string, handler httprouter.Handle) {
	if c.routes == nil {
		c.routes = map[string]map[string]httprouter.Handle{}
	}
	if c.routes[method] == nil {
		c.routes[method] = map[string]httprouter.Handle{}
	}
	c.routes[method][path] = handler
}

func (c *HttpServerConfig) GetRoutes() map[string]map[string]httprouter.Handle {
	return c.routes
}
