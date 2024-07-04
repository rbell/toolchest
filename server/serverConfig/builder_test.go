/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package serverConfig

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/rbell/toolchest/server/example/grpcService"
	"github.com/rbell/toolchest/server/example/proto"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestServerConfigBuilder_Build_BuildsHttpServiceConfig(t *testing.T) {
	// setup

	// test
	cfg := BuildServerConfig().WithHttpServiceConfig(
		BuildHttpServiceConfig().
			WithPort("8080").
			AddRoute("GET", "/hello", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
				//nolint:errcheck // ignore lint error for example
				w.Write([]byte("Hello, World!"))
			})).Build().GetHttpServerConfig()

	// verify
	assert.NotNil(t, cfg)
	assert.Equal(t, "8080", cfg.Port)
	assert.NotNil(t, cfg.GetRoutes())
	assert.Len(t, cfg.GetRoutes(), 1)
}

func TestServerConfigBuilder_Build_BuildsHttpsServiceConfig(t *testing.T) {
	// setup
	tlsCfg := &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	// test
	cfg := BuildServerConfig().WithHttpsServiceConfig(
		BuildHttpsServiceConfig().
			WithPort("8443").
			AddRoute("GET", "/hello", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
				//nolint:errcheck // ignore lint error for example
				w.Write([]byte("Hello, World!"))
			}).
			WithTlsConfig(tlsCfg).
			WithCertFile("testCert").
			WithKeyFile("testKey")).Build().GetHttpsServerConfig()

	// verify
	assert.NotNil(t, cfg)
	assert.Equal(t, "8443", cfg.Port)
	assert.NotNil(t, cfg.GetRoutes())
	assert.Len(t, cfg.GetRoutes(), 1)
	assert.Equal(t, tlsCfg, cfg.GetTlsConfig())
	assert.Equal(t, "testCert", cfg.GetCertFile())
	assert.Equal(t, "testKey", cfg.GetKeyFile())
}

func TestServerConfigBuilder_Build_BuildsGrpcServiceConfig(t *testing.T) {
	// setup

	// test
	cfg := BuildServerConfig().WithGrpcServiceConfig(
		BuildGrpcServerConfig().
			WithPort("8888").
			RegisterImplementation(&proto.HelloService_ServiceDesc, &grpcService.HelloService{}).
			EnableReflection()).Build().GetGrpcServerConfig()

	// verify
	assert.NotNil(t, cfg)
	assert.Equal(t, "8888", cfg.Port)
	assert.True(t, cfg.IsReflectionEnabled())
	assert.NotEmpty(t, cfg.GetRegistrations())
}
