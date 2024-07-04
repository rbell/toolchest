/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package grpcServer

import (
	"context"
	"log/slog"
	"net"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

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
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// test
	provider.Start(context.Background(), wg, wg)

	// verify
	mGrpcServer.AssertExpectations(t)
}

func TestGrpcProvider_Stop(t *testing.T) {
	// setup
	mGrpcServer := &mockGrpcServerer{}
	mGrpcServer.On("GracefulStop").Return()
	provider := &GrpcProvider{
		grpcServer: mGrpcServer,
		logger:     slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	// test
	err := provider.Stop(context.Background())

	// verify
	assert.Nil(t, err)
	mGrpcServer.AssertExpectations(t)
}
