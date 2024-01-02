/*
 * Copyright (c) 2024  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package propositions

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapContainsKey_ReturnsTrue_WhenKeyExists(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	result := MapContainsKey(m, 2)
	assert.True(t, result)
}

func TestMapContainsKey_ReturnsFalse_WhenKeyDoesNotExist(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	result := MapContainsKey(m, 4)
	assert.False(t, result)
}

func TestMapContainsValue_ReturnsTrue_WhenValueExists(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	result := MapContainsValue(m, "two")
	assert.True(t, result)
}

func TestMapContainsValue_ReturnsFalse_WhenValueDoesNotExist(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	result := MapContainsValue(m, "four")
	assert.False(t, result)
}

func TestMapKeyPropositionAny_ReturnsTrue_WhenAnyKeySatisfiesProposition(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	result := MapKeyPropositionAny(m, func(e int) bool {
		return e > 2
	})
	assert.True(t, result)
}

func TestMapKeyPropositionAny_ReturnsFalse_WhenNoKeySatisfiesProposition(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	result := MapKeyPropositionAny(m, func(e int) bool {
		return e > 3
	})
	assert.False(t, result)
}

func TestMapKeyPropositionAll_ReturnsTrue_WhenAllKeysSatisfyProposition(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	result := MapKeyPropositionAll(m, func(e int) bool {
		return e < 4
	})
	assert.True(t, result)
}

func TestMapKeyPropositionAll_ReturnsFalse_WhenAnyKeyDoesNotSatisfyProposition(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	result := MapKeyPropositionAll(m, func(e int) bool {
		return e < 3
	})
	assert.False(t, result)
}

func TestMapKeyPropositionNone_ReturnsTrue_WhenNoKeySatisfiesProposition(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	result := MapKeyPropositionNone(m, func(e int) bool {
		return e > 3
	})
	assert.True(t, result)
}

func TestMapKeyPropositionNone_ReturnsFalse_WhenAnyKeySatisfiesProposition(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	result := MapKeyPropositionNone(m, func(e int) bool {
		return e > 2
	})
	assert.False(t, result)
}

func TestMapValuePropositionAny_ReturnsTrue_WhenAnyValueSatisfiesProposition(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	result := MapValuePropositionAny(m, func(e string) bool {
		return e == "three"
	})
	assert.True(t, result)
}

func TestMapValuePropositionAny_ReturnsFalse_WhenNoValueSatisfiesProposition(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	result := MapValuePropositionAny(m, func(e string) bool {
		return e == "four"
	})
	assert.False(t, result)
}

func TestMapValuePropositionAll_ReturnsTrue_WhenAllValuesSatisfyProposition(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	result := MapValuePropositionAll(m, func(e string) bool {
		return len(e) > 0
	})
	assert.True(t, result)
}

func TestMapValuePropositionAll_ReturnsFalse_WhenAnyValueDoesNotSatisfyProposition(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	result := MapValuePropositionAll(m, func(e string) bool {
		return e == "one"
	})
	assert.False(t, result)
}

func TestMapValuePropositionNone_ReturnsTrue_WhenNoValueSatisfiesProposition(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	result := MapValuePropositionNone(m, func(e string) bool {
		return e == "four"
	})
	assert.True(t, result)
}

func TestMapValuePropositionNone_ReturnsFalse_WhenAnyValueSatisfiesProposition(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	result := MapValuePropositionNone(m, func(e string) bool {
		return e == "three"
	})
	assert.False(t, result)
}
