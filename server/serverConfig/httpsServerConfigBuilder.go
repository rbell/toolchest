/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package serverConfig

import (
	"crypto/tls"
	"net/http"
)

type HttpsServerConfigBuilder struct {
	cfg *HttpsServerConfig
}

func BuildHttpsServiceConfig() *HttpsServerConfigBuilder {
	return &HttpsServerConfigBuilder{cfg: &HttpsServerConfig{
		HttpServerConfig: BuildHttpServiceConfig().build(),
	}}
}

func (b *HttpsServerConfigBuilder) WithPort(port string) *HttpsServerConfigBuilder {
	b.cfg.Port = port
	return b
}

func (b *HttpsServerConfigBuilder) AddRoute(method, path string, handler http.HandlerFunc) *HttpsServerConfigBuilder {
	b.cfg.AddRoute(method, path, handler)
	return b
}

func (b *HttpsServerConfigBuilder) WithTlsConfig(tlsConfig *tls.Config) *HttpsServerConfigBuilder {
	b.cfg.SetTlsConfig(tlsConfig)
	return b
}

func (b *HttpsServerConfigBuilder) WithCertFile(certFile string) *HttpsServerConfigBuilder {
	b.cfg.SetCertFile(certFile)
	return b
}

func (b *HttpsServerConfigBuilder) WithKeyFile(keyFile string) *HttpsServerConfigBuilder {
	b.cfg.SetKeyFile(keyFile)
	return b
}

func (b *HttpsServerConfigBuilder) build() *HttpsServerConfig {
	return b.cfg
}
