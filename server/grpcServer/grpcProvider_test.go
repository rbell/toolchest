/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package grpcServer

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/rbell/toolchest/server/grpcServer/netMocks"
	"github.com/rbell/toolchest/server/serverConfig"
)

func TestGrpcProvider_Start(t *testing.T) {
	// setup
	port := "8080"
	mListener := &netMocks.Listener{}

	mGrpcServer := &mockGrpcServerer{}
	mGrpcServer.On("Serve", mListener).Return(nil)

	provider := &GrpcProvider{
		grpcServer: mGrpcServer,
		cfg:        &serverConfig.GrpcServerConfig{Port: port},
		listenerFactory: func(string, string) (net.Listener, error) {
			return mListener, nil
		},
	}
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// test
	provider.Start(wg, wg)

	// verify
	mGrpcServer.AssertExpectations(t)
}

func TestGrpcProvider_Stop(t *testing.T) {
	// setup
	mGrpcServer := &mockGrpcServerer{}
	mGrpcServer.On("GracefulStop").Return()
	provider := &GrpcProvider{
		grpcServer: mGrpcServer,
	}
	// test
	provider.Stop(context.Background())

	// verify
	mGrpcServer.AssertExpectations(t)
}
