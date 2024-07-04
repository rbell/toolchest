/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package httpServer

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpProvider_Start_NonTLS_ListenAndServes(t *testing.T) {
	// setup
	wg := &sync.WaitGroup{}
	wg.Add(2)
	mListener := &mockHttpListener{}
	mListener.On("ListenAndServe").Return(nil)

	provider := &HttpProvider{
		httpSrver: mListener,
		isTLS:     false,
	}

	// test
	provider.Start(wg, wg)

	// verify
	mListener.AssertExpectations(t)
}

func TestHttpProvider_Start_TLS_ListenAndServes(t *testing.T) {
	// setup
	wg := &sync.WaitGroup{}
	wg.Add(2)
	mListener := &mockHttpListener{}
	mListener.On("ListenAndServeTLS", "", "").Return(nil)

	provider := &HttpProvider{
		httpSrver: mListener,
		isTLS:     true,
	}

	// test
	provider.Start(wg, wg)

	// verify
	mListener.AssertExpectations(t)
}

func TestHttpProvider_Stop_ClosesServer(t *testing.T) {
	// setup
	wg := &sync.WaitGroup{}
	wg.Add(2)
	ctx := context.Background()
	mListener := &mockHttpListener{}
	mListener.On("Shutdown", ctx).Return(nil)

	provider := &HttpProvider{
		httpSrver: mListener,
	}

	// test
	err := provider.Stop(ctx)

	// verify
	assert.Nil(t, err)
	mListener.AssertExpectations(t)
}
