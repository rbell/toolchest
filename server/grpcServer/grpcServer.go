/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package server

import (
	"context"
	"github.com/rbell/toolchest/server/serverConfig"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type GrpcProvider struct {
	grpcServer *grpc.Server
	cfg        *serverConfig.GrpcServerConfig
}

func NewGrpcProvider(cfg *serverConfig.GrpcServerConfig) *GrpcProvider {
	return &GrpcProvider{
		grpcServer: grpc.NewServer(cfg.GetOpts()...),
		cfg:        cfg,
	}
}

func (p *GrpcProvider) Start(startWg, stopWg *sync.WaitGroup) {
	defer stopWg.Done()
	ls, err := net.Listen("tcp", p.cfg.Port)
	if err != nil {
		log.Fatal(err)
	}

	startWg.Done()
	err = p.grpcServer.Serve(ls)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *GrpcProvider) Stop(ctx context.Context) {
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
}
