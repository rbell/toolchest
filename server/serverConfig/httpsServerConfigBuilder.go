/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package serverConfig

import (
	"crypto/tls"
	"github.com/julienschmidt/httprouter"
)

type HttpsServerConfigBuilder struct {
	cfg *HttpsServerConfig
}

func BuildHttpsServiceConfig() *HttpsServerConfigBuilder {
	return &HttpsServerConfigBuilder{cfg: &HttpsServerConfig{}}
}

func (b *HttpsServerConfigBuilder) WithPort(port string) *HttpsServerConfigBuilder {
	b.cfg.Port = port
	return b
}

func (b *HttpsServerConfigBuilder) AddRoute(method, path string, handler httprouter.Handle) *HttpsServerConfigBuilder {
	if b.cfg.routes == nil {
		b.cfg.routes = map[string]map[string]httprouter.Handle{}
	}
	if b.cfg.routes[method] == nil {
		b.cfg.routes[method] = map[string]httprouter.Handle{}
	}
	b.cfg.routes[method][path] = handler
	return b
}

func (b *HttpsServerConfigBuilder) WithTlsConfig(tlsConfig *tls.Config) *HttpsServerConfigBuilder {
	b.cfg.tlsConfig = tlsConfig
	return b
}

func (b *HttpsServerConfigBuilder) WithCertFile(certFile string) *HttpsServerConfigBuilder {
	b.cfg.certFile = certFile
	return b
}

func (b *HttpsServerConfigBuilder) WithKeyFile(keyFile string) *HttpsServerConfigBuilder {
	b.cfg.keyFile = keyFile
	return b
}

func (b *HttpsServerConfigBuilder) build() *HttpsServerConfig {
	return b.cfg
}
