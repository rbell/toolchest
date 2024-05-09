/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package objectSearch

import (
	"cmp"
	"reflect"
)

type trait struct {
	name      string
	traitType reflect.Type
}

func newTrait[T cmp.Ordered](name string) *trait {
	var t T
	tt := reflect.TypeOf(t)
	return &trait{
		name: name,
		traitType: tt,
	}
}

func (t *trait[T])

type index struct {
}
