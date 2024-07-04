/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package server

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"github.com/rbell/toolchest/server/grpcServer"
	"github.com/rbell/toolchest/server/httpServer"
	"github.com/rbell/toolchest/server/serverConfig"
	"github.com/richardwilkes/toolbox/errs"
)

type ServiceProvider interface {
	Start(ctx context.Context, startWg, stopWg *sync.WaitGroup)
	Stop(ctx context.Context) error
}

type Server struct {
	providers       []ServiceProvider
	stopProvidersWg *sync.WaitGroup
	runningCtx      context.Context
	logger          *slog.Logger
}

func NewServer(cfg *serverConfig.Config, runningCtx context.Context, stopWg *sync.WaitGroup) (*Server, error) {
	logger := cfg.GetLogger()
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	providers := []ServiceProvider{}

	if httpCfg := cfg.GetHttpServerConfig(); httpCfg != nil {
		providers = append(providers, httpServer.NewHttpProvider(httpCfg, logger))
	}

	if httpsCfg := cfg.GetHttpsServerConfig(); httpsCfg != nil {
		providers = append(providers, httpServer.NewHttpsProvider(httpsCfg, logger))
	}

	if grpcCfg := cfg.GetGrpcServerConfig(); grpcCfg != nil {
		providers = append(providers, grpcServer.NewGrpcProvider(grpcCfg, logger))
	}

	return &Server{
		providers:       providers,
		stopProvidersWg: stopWg,
		runningCtx:      runningCtx,
		logger:          logger,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	startWg := &sync.WaitGroup{}
	for _, provider := range s.providers {
		s.stopProvidersWg.Add(1)
		startWg.Add(1)
		go provider.Start(ctx, startWg, s.stopProvidersWg)
	}

	startWg.Wait()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	var err error
	for _, provider := range s.providers {
		e := provider.Stop(ctx)
		if e != nil {
			err = errs.Append(err, e)
		}
	}
	s.stopProvidersWg.Wait()
	return err
}
