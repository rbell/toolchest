/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package httpServer

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rbell/toolchest/server/serverConfig"
	"log"
	"net/http"
	"sync"
)

type HttpProvider struct {
	httpSrver *http.Server
	router    *httprouter.Router
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

func (p *HttpProvider) Start(startWg, stopWg *sync.WaitGroup) {
	defer stopWg.Done()

	startWg.Done()
	err := p.httpSrver.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (p *HttpProvider) Stop(ctx context.Context) {
	p.httpSrver.Shutdown(ctx)
}
