/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package serverConfig

type Config struct {
	httpServiceConfig *HttpServerConfig
	grpcServerConfig  *GrpcServerConfig
}

func (c *Config) Validate() error {
	return nil
}

func (c *Config) GetHttpServiceConfig() *HttpServerConfig {
	return c.httpServiceConfig
}

func (c *Config) GetGrpcServerConfig() *GrpcServerConfig {
	return c.grpcServerConfig
}
