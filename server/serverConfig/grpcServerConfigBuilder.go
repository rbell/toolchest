/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package serverConfig

import "google.golang.org/grpc"

type GrpcServerConfigBuilder struct {
	cfg *GrpcServerConfig
}

func BuildGrpcServerConfig() *GrpcServerConfigBuilder {
	return &GrpcServerConfigBuilder{cfg: &GrpcServerConfig{}}
}

func (b *GrpcServerConfigBuilder) WithPort(port string) *GrpcServerConfigBuilder {
	b.cfg.Port = port
	return b
}

func (b *GrpcServerConfigBuilder) AddOption(opt grpc.ServerOption) *GrpcServerConfigBuilder {
	b.cfg.AddOption(opt)
	return b
}

func (b *GrpcServerConfigBuilder) RegisterImplementation(description *grpc.ServiceDesc, impl any) *GrpcServerConfigBuilder {
	b.cfg.RegisterImplementation(description, impl)
	return b
}

func (b *GrpcServerConfigBuilder) AddInitializer(init GrpcServerInitializer) *GrpcServerConfigBuilder {
	b.cfg.AddInitializer(init)
	return b
}

func (b *GrpcServerConfigBuilder) build() *GrpcServerConfig {
	return b.cfg
}

func (b *GrpcServerConfigBuilder) EnableReflection() *GrpcServerConfigBuilder {
	b.cfg.EnableReflection()
	return b
}
