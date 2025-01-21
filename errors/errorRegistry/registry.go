/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package errorRegistry

import (
	"context"

	"github.com/rbell/toolchest/stacktrace"
)

func init() {
	reg = newRegistry()
}

var reg *registry

type ErrorResolver interface {
	ResolveError(error, context.Context, stacktrace.StackTrace) error
}

type registry struct {
	resolvers       []ErrorResolver
	defaultResolver ErrorResolver
}

func newRegistry() *registry {
	return &registry{
		resolvers: make([]ErrorResolver, 0),
	}
}

func RegisterResolver(resolver ErrorResolver) {
	reg.resolvers = append(reg.resolvers, resolver)
}

func DefaultResolver(resolver ErrorResolver) {
	reg.defaultResolver = resolver
}

func ResolveError(err error, ctx context.Context) error {
	st := stacktrace.CaptureStackTrace()
	for _, resolver := range reg.resolvers {
		if resolved := resolver.ResolveError(err, ctx, st); resolved != nil {
			return resolved
		}
	}
	if reg.defaultResolver != nil {
		return reg.defaultResolver.ResolveError(err, ctx, st)
	}
	return err
}

func ClearResolvers() {
	reg.resolvers = make([]ErrorResolver, 0)
}
