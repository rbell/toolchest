/*
 * Copyright (c) 2025 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package errorRegistry

import (
	"context"
	"errors"

	"github.com/rbell/toolchest/stacktrace"
)

type errorConstructor[S error, T error] func(S, context.Context) T

type GenericResolver[S error, T error] struct {
	constructor    errorConstructor[S, T]
	stackPredicate func(st stacktrace.StackTrace) bool
}

type Option[S error, T error] func(*GenericResolver[S, T])

func NewGenericResolver[S error, T error](constructor errorConstructor[S, T], opts ...Option[S, T]) *GenericResolver[S, T] {
	resolver := &GenericResolver[S, T]{
		constructor: constructor,
	}
	for _, opt := range opts {
		opt(resolver)
	}
	return resolver
}

func (r *GenericResolver[S, T]) ResolveError(err error, ctx context.Context, st stacktrace.StackTrace) error {
	var s S
	if errors.As(err, &s) {
		return r.constructor(err.(S), ctx)
	}
	return nil
}

func (r *GenericResolver[S, T]) WhenStackReferencesFile(fileEnding string) Option[S, T] {
	return func(resolver *GenericResolver[S, T]) {
		resolver.stackPredicate = func(st stacktrace.StackTrace) bool {
			return st.ReferencesFile(fileEnding)
		}
	}
}

func (r *GenericResolver[S, T]) WhenStackReferencesFunction(funcName string) Option[S, T] {
	return func(resolver *GenericResolver[S, T]) {
		resolver.stackPredicate = func(st stacktrace.StackTrace) bool {
			return st.ReferencesFunction(funcName)
		}
	}
}
