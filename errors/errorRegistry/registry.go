/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package errorRegistry

func init() {
	reg = newRegistry()
}

var reg *registry

type ErrorResolver interface {
	ResolveError(error error) error
}

type registry struct {
	resolvers []ErrorResolver
}

func newRegistry() *registry {
	return &registry{
		resolvers: make([]ErrorResolver, 0),
	}
}

func RegisterResolver(resolver ErrorResolver) {
	reg.resolvers = append(reg.resolvers, resolver)
}
