/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package serverConfig

import "google.golang.org/grpc"

type GrpcServerInitializer func(server *grpc.Server)

type GrpcServerConfig struct {
	Port         string
	opts         []grpc.ServerOption
	initializers []GrpcServerInitializer
}

func (c *GrpcServerConfig) AddOption(opt grpc.ServerOption) {
	c.opts = append(c.opts, opt)
}

func (c *GrpcServerConfig) GetOpts() []grpc.ServerOption {
	return c.opts
}

func (c *GrpcServerConfig) AddInitializer(init GrpcServerInitializer) {
	c.initializers = append(c.initializers, init)
}

func (c *GrpcServerConfig) GetInitializers() []GrpcServerInitializer {
	return c.initializers
}
