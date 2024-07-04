/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package grpcServer

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/rbell/toolchest/server/serverConfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServerer interface {
	Serve(lis net.Listener) error
	Stop()
	GracefulStop()
}

type listenerFactory func(string, string) (net.Listener, error)

type GrpcProvider struct {
	grpcServer      grpcServerer
	cfg             *serverConfig.GrpcServerConfig
	listenerFactory listenerFactory
}

func NewGrpcProvider(cfg *serverConfig.GrpcServerConfig) *GrpcProvider {
	srvr := grpc.NewServer(cfg.GetOpts()...)
	provider := &GrpcProvider{
		grpcServer:      srvr,
		cfg:             cfg,
		listenerFactory: net.Listen,
	}

	for desc, impl := range cfg.GetRegistrations() {
		srvr.RegisterService(desc, impl)
	}

	if cfg.IsReflectionEnabled() {
		reflection.Register(srvr)
	}

	for _, initializer := range cfg.GetInitializers() {
		initializer(srvr)
	}
	return provider
}

func (p *GrpcProvider) Start(startWg, stopWg *sync.WaitGroup) {
	defer stopWg.Done()
	ls, err := p.listenerFactory("tcp", fmt.Sprintf(":%s", p.cfg.Port))
	if err != nil {
		log.Fatal(err)
	}

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
