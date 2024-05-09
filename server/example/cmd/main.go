/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package main

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/rbell/toolchest/server"
	"github.com/rbell/toolchest/server/serverConfig"
	"github.com/richardwilkes/toolbox/atexit"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	cfg := serverConfig.BuildServerConfig().
		WithHttpServiceConfig(
			serverConfig.BuildHttpServiceConfig().
				WithPort("8080").
				AddRoute("GET", "/hello", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
					w.Write([]byte("Hello, World!"))
				})).
		Build()

	wg := &sync.WaitGroup{}
	srvr, err := server.NewServer(cfg, context.Background(), wg)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	atexit.Register(func() {
		stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Second*30)
		defer stopCancel()
		srvr.Stop(stopCtx)
	})

	srvr.Start()

	wg.Wait()
}
