/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package httpServer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/rbell/toolchest/server/serverConfig"
)

type httpServerOption func(server *http.Server)

type HttpProvider struct {
	httpSrver *http.Server
	router    *httprouter.Router
	certFile  string
	keyFile   string
}

func NewHttpProvider(cfg *serverConfig.HttpServerConfig) *HttpProvider {
	provider := &HttpProvider{
		httpSrver: &http.Server{Addr: fmt.Sprintf(":%v", cfg.Port)},
		router:    httprouter.New(),
	}

	for method, paths := range cfg.GetRoutes() {
		for path, handler := range paths {
			provider.router.Handle(method, path, handler)
		}
	}

	provider.httpSrver.Handler = provider.router

	return provider
}

func NewHttpsProvider(cfg *serverConfig.HttpsServerConfig) *HttpProvider {
	provider := NewHttpProvider(cfg.HttpServerConfig)
	provider.httpSrver.TLSConfig = cfg.GetTlsConfig()
	provider.certFile = cfg.GetCertFile()
	provider.keyFile = cfg.GetKeyFile()

	return provider
}

func (p *HttpProvider) Start(startWg, stopWg *sync.WaitGroup) {
	defer stopWg.Done()

	startWg.Done()
	fmt.Println("HTTP Server started on port: ", p.httpSrver.Addr)

	if p.httpSrver.TLSConfig != nil {
		err := p.httpSrver.ListenAndServeTLS(p.certFile, p.keyFile)
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return
	}

	err := p.httpSrver.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (p *HttpProvider) Stop(ctx context.Context) error {
	return p.httpSrver.Shutdown(ctx)
}
