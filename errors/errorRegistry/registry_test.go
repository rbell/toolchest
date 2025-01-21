/*
 * Copyright (c) 2025 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package errorRegistry

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/richardwilkes/toolbox/errs"
)

type embeddedErr = errs.Error

type TestErrorOne struct {
	*embeddedErr
}

type TestErrorTwo struct {
	*embeddedErr
}

func TestRegisterResolver_RegistersResolver(t *testing.T) {
	resolver := NewGenericResolver(func(err *TestErrorOne, ctx context.Context) *TestErrorTwo {
		return &TestErrorTwo{
			embeddedErr: errs.NewWithCause("TestErrorTwo", err),
		}
	})

	RegisterResolver(resolver)

	err := &TestErrorOne{errs.NewWithCause("TestErrorOne", nil)}
	resolved := ResolveError(err, context.Background())
	assert.True(t, strings.HasPrefix(resolved.Error(), "TestErrorTwo"))
}
