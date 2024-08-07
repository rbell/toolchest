/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/rbell/toolchest/server/httpMiddleware"

	"github.com/rbell/toolchest/server"
	"github.com/rbell/toolchest/server/example/grpcService"
	"github.com/rbell/toolchest/server/example/proto"
	"github.com/rbell/toolchest/server/serverConfig"
	"github.com/richardwilkes/toolbox/atexit"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// set up middleware to log both request and response.
	// Don't like the way the middleware logs? Implement your own middleware using the ones supplied as patterns.
	middleware := httpMiddleware.BundleMiddleware(httpMiddleware.LogRequest(logger), httpMiddleware.LogResponse(logger))

	cfg := serverConfig.BuildServerConfig().
		WithLogger(logger).
		WithHttpServiceConfig(
			serverConfig.BuildHttpServiceConfig().
				UsingMiddleWare(middleware).
				WithPort("8080").
				AddRoute("GET", "/hello", func(w http.ResponseWriter, r *http.Request) {
					//nolint:errcheck // ignore lint error for example
					w.Write([]byte("Hello, World!"))
				})).
		WithGrpcServiceConfig(
			serverConfig.BuildGrpcServerConfig().
				WithPort("8888").
				RegisterImplementation(&proto.HelloService_ServiceDesc, &grpcService.HelloService{}).
				EnableReflection()).
		Build()

	// wait group to wait for server to stop
	wg := &sync.WaitGroup{}
	// create server, passing in configuration, startup context and wait group.
	srvr, err := server.NewServer(cfg, context.Background(), wg)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	// register a function to stop the server when the application exits (in this example we use the atexit library to register the function)
	atexit.Register(func() {
		// context with timeout for server to stop
		stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Second*30)
		//nolint:errcheck // ignore lint error for example
		defer stopCancel()

		//nolint:errcheck // ignore errorlint error for example
		srvr.Stop(stopCtx)
	})

	// start the server
	err = srvr.Start(context.Background())
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	// wait for server to stop
	wg.Wait()
}
