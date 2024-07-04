/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package httpServer

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/rbell/toolchest/server/serverConfig"
)

type httpListener interface {
	ListenAndServe() error
	ListenAndServeTLS(certFile, keyFile string) error
	Shutdown(ctx context.Context) error
}

type HttpProvider struct {
	httpSrver httpListener
	address   string
	isTLS     bool
	certFile  string
	keyFile   string
	logger    *slog.Logger
}

func NewHttpProvider(cfg *serverConfig.HttpServerConfig, logger *slog.Logger) *HttpProvider {
	srvr := &http.Server{Addr: fmt.Sprintf(":%v", cfg.Port)}
	provider := &HttpProvider{
		httpSrver: srvr,
		isTLS:     false,
		address:   fmt.Sprintf(":%v", cfg.Port),
		logger:    logger,
	}

	router := httprouter.New()

	for method, paths := range cfg.GetRoutes() {
		for path, handler := range paths {
			router.Handle(method, path, handler)
		}
	}

	srvr.Handler = router

	return provider
}

func NewHttpsProvider(cfg *serverConfig.HttpsServerConfig, logger *slog.Logger) *HttpProvider {
	srvr := &http.Server{Addr: fmt.Sprintf(":%v", cfg.Port)}
	srvr.TLSConfig = cfg.GetTlsConfig()

	provider := &HttpProvider{
		httpSrver: srvr,
		certFile:  cfg.GetCertFile(),
		keyFile:   cfg.GetKeyFile(),
		isTLS:     true,
		address:   fmt.Sprintf(":%v", cfg.Port),
		logger:    logger,
	}

	router := httprouter.New()

	for method, paths := range cfg.GetRoutes() {
		for path, handler := range paths {
			router.Handle(method, path, handler)
		}
	}

	return provider
}

func (p *HttpProvider) Start(ctx context.Context, startWg, stopWg *sync.WaitGroup) {
	defer stopWg.Done()

	startWg.Done()

	if p.isTLS {
		p.logger.InfoContext(ctx, fmt.Sprintf("Starting HTTPS server on %v", p.address))
		err := p.httpSrver.ListenAndServeTLS(p.certFile, p.keyFile)
		if err != nil && err != http.ErrServerClosed {
			p.logger.ErrorContext(ctx, fmt.Sprintf("Error Starting HTTPS server on %v", p.address), slog.Any("error", err.Error()))
		}
		return
	}

	p.logger.InfoContext(ctx, fmt.Sprintf("Starting HTTP server on %v", p.address))
	err := p.httpSrver.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		p.logger.ErrorContext(ctx, fmt.Sprintf("Error Starting HTTP server on %v", p.address), slog.Any("error", err.Error()))
	}
}

func (p *HttpProvider) Stop(ctx context.Context) error {
	defer func() {
		if p.isTLS {
			p.logger.InfoContext(ctx, "Stopped HTTPS server")
		} else {
			p.logger.InfoContext(ctx, "Stopped HTTP server")
		}
	}()

	return p.httpSrver.Shutdown(ctx)
}
