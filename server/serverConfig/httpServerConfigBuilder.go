/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package serverConfig

import "github.com/julienschmidt/httprouter"

type HttpServerConfigBuilder struct {
	cfg *HttpServerConfig
}

func BuildHttpServiceConfig() *HttpServerConfigBuilder {
	return &HttpServerConfigBuilder{cfg: &HttpServerConfig{}}
}

func (b *HttpServerConfigBuilder) WithPort(port string) *HttpServerConfigBuilder {
	b.cfg.Port = port
	return b
}

func (b *HttpServerConfigBuilder) AddRoute(method, path string, handler httprouter.Handle) *HttpServerConfigBuilder {
	if b.cfg.routes == nil {
		b.cfg.routes = map[string]map[string]httprouter.Handle{}
	}
	if b.cfg.routes[method] == nil {
		b.cfg.routes[method] = map[string]httprouter.Handle{}
	}
	b.cfg.routes[method][path] = handler
	return b
}

func (b *HttpServerConfigBuilder) build() *HttpServerConfig {
	return b.cfg
}
