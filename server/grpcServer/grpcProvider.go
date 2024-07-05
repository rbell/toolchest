/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package grpcServer

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"github.com/rbell/toolchest/server/internal/sharedTypes"

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
	logger          sharedTypes.LogPublisher
}

func NewGrpcProvider(cfg *serverConfig.GrpcServerConfig, logger sharedTypes.LogPublisher) *GrpcProvider {
	srvr := grpc.NewServer(cfg.GetOpts()...)
	provider := &GrpcProvider{
		grpcServer:      srvr,
		cfg:             cfg,
		listenerFactory: net.Listen,
		logger:          logger,
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

func (p *GrpcProvider) Start(ctx context.Context, startWg, stopWg *sync.WaitGroup) {
	defer stopWg.Done()

	p.logger.InfoContext(ctx, fmt.Sprintf("Starting gRPC server on %v", p.cfg.Port))
	ls, err := p.listenerFactory("tcp", fmt.Sprintf(":%s", p.cfg.Port))
	if err != nil {
		p.logger.ErrorContext(ctx, "Error Starting Listener for gRPC server", slog.Any("error", err.Error()))
	}

	startWg.Done()
	err = p.grpcServer.Serve(ls)
	if err != nil {
		p.logger.ErrorContext(ctx, "Error Starting gRPC server", slog.Any("error", err.Error()))
	}
}

func (p *GrpcProvider) Stop(ctx context.Context) error {
	stoppedCh := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			p.logger.WarnContext(ctx, "Forcefully Stopping gRPC server")
			p.grpcServer.Stop()
		case <-stoppedCh:
		}
	}()
	p.grpcServer.GracefulStop()
	p.logger.InfoContext(ctx, "Stopped gRPC server")
	close(stoppedCh)
	return nil
}
