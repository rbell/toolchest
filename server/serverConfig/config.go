/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package serverConfig

import "log/slog"

type Config struct {
	httpServerConfig  *HttpServerConfig
	httpsServerConfig *HttpsServerConfig
	grpcServerConfig  *GrpcServerConfig
	logger            *slog.Logger
}

func (c *Config) Validate() error {
	return nil
}

func (c *Config) GetHttpServerConfig() *HttpServerConfig {
	return c.httpServerConfig
}

func (c *Config) GetHttpsServerConfig() *HttpsServerConfig {
	return c.httpsServerConfig
}

func (c *Config) GetGrpcServerConfig() *GrpcServerConfig {
	return c.grpcServerConfig
}

func (c *Config) GetLogger() *slog.Logger {
	return c.logger
}
