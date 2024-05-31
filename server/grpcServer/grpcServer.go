/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package grpcServer

import (
	"context"
	"fmt"
	"github.com/rbell/toolchest/server/serverConfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
)

type GrpcProvider struct {
	grpcServer *grpc.Server
	cfg        *serverConfig.GrpcServerConfig
}

func NewGrpcProvider(cfg *serverConfig.GrpcServerConfig) *GrpcProvider {
	provider := &GrpcProvider{
		grpcServer: grpc.NewServer(cfg.GetOpts()...),
		cfg:        cfg,
	}

	for desc, impl := range cfg.GetRegistrations() {
		provider.grpcServer.RegisterService(desc, impl)
	}

	if cfg.IsReflectionEnabled() {
		reflection.Register(provider.grpcServer)
	}

	for _, initializer := range cfg.GetInitializers() {
		initializer(provider.grpcServer)
	}
	return provider
}

func (p *GrpcProvider) Start(startWg, stopWg *sync.WaitGroup) {
	defer stopWg.Done()
	ls, err := net.Listen("tcp", fmt.Sprintf(":%s", p.cfg.Port))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("GRPC Server started on port: ", p.cfg.Port)
	startWg.Done()
	err = p.grpcServer.Serve(ls)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *GrpcProvider) Stop(ctx context.Context) error {
	stoppedCh := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			p.grpcServer.Stop()
		case <-stoppedCh:
		}
	}()
	p.grpcServer.GracefulStop()
	close(stoppedCh)
	return nil
}
