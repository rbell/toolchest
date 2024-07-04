/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package serverConfig

type ServerConfigBuilder struct {
	cfg *Config
}

func BuildServerConfig() *ServerConfigBuilder {
	return &ServerConfigBuilder{cfg: &Config{}}
}

func (b *ServerConfigBuilder) WithHttpServiceConfig(httpBuilder *HttpServerConfigBuilder) *ServerConfigBuilder {
	b.cfg.httpServerConfig = httpBuilder.build()
	return b
}

func (b *ServerConfigBuilder) WithHttpsServiceConfig(httpsBuilder *HttpsServerConfigBuilder) *ServerConfigBuilder {
	b.cfg.httpsServerConfig = httpsBuilder.build()
	return b
}

func (b *ServerConfigBuilder) WithGrpcServiceConfig(grpcBuilder *GrpcServerConfigBuilder) *ServerConfigBuilder {
	b.cfg.grpcServerConfig = grpcBuilder.build()
	return b
}

func (b *ServerConfigBuilder) Build() *Config {
	return b.cfg
}
