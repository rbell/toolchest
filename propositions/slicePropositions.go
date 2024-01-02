/*
 * Copyright (c) 2023-2024  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package propositions

import (
	"cmp"
)

// SliceContains returns true if the slice contains the value.
func SliceContains[T cmp.Ordered](s []T, v T) bool {
	return SlicePropositionAny(s, func(e T) bool {
		return e == v
	})
}

// SliceContainsAll returns true if the slice contains all the values.
func SliceContainsAll[T cmp.Ordered](s []T, v []T) bool {
	for _, e := range v {
		if !SliceContains(s, e) {
			return false
		}
	}
	return true
}

// SliceContainsAny returns true if the slice contains any of the values.
func SliceContainsAny[T cmp.Ordered](s []T, v []T) bool {
	for _, e := range v {
		if SliceContains(s, e) {
			return true
		}
	}
	return false
}

// SliceContainsNone returns true if the slice contains none of the values.
func SliceContainsNone[T cmp.Ordered](s []T, v []T) bool {
	for _, e := range v {
		if SliceContains(s, e) {
			return false
		}
	}
	return true
}

// SliceAllLessThan returns true if all the values in the slice are less than the value.
func SliceAllLessThan[T cmp.Ordered](s []T, v T) bool {
	return SlicePropositionAll(s, func(e T) bool {
		return e < v
	})
}

// SliceAllLessThanOrEqualTo returns true if all the values in the slice are less than or equal to the value.
func SliceAllLessThanOrEqualTo[T cmp.Ordered](s []T, v T) bool {
	return SlicePropositionAll(s, func(e T) bool {
		return e <= v
	})
}

// SliceAllGreaterThan returns true if all the values in the slice are greater than the value.
func SliceAllGreaterThan[T cmp.Ordered](s []T, v T) bool {
	return SlicePropositionAll(s, func(e T) bool {
		return e > v
	})
}

// SliceAllGreaterThanOrEqualTo returns true if all the values in the slice are greater than or equal to the value.
func SliceAllGreaterThanOrEqualTo[T cmp.Ordered](s []T, v T) bool {
	return SlicePropositionAll(s, func(e T) bool {
		return e >= v
	})
}

// SliceAnyLessThan returns true if any of the values in the slice are less than the value.
func SliceAnyLessThan[T cmp.Ordered](s []T, v T) bool {
	return SlicePropositionAny(s, func(e T) bool {
		return e < v
	})
}

// SliceAnyLessThanOrEqualTo returns true if any of the values in the slice are less than or equal to the value.
func SliceAnyLessThanOrEqualTo[T cmp.Ordered](s []T, v T) bool {
	return SlicePropositionAny(s, func(e T) bool {
		return e <= v
	})
}

// SliceAnyGreaterThan returns true if any of the values in the slice are greater than the value.
func SliceAnyGreaterThan[T cmp.Ordered](s []T, v T) bool {
	return SlicePropositionAny(s, func(e T) bool {
		return e > v
	})
}

// SliceAnyGreaterThanOrEqualTo returns true if any of the values in the slice are greater than or equal to the value.
func SliceAnyGreaterThanOrEqualTo[T cmp.Ordered](s []T, v T) bool {
	return SlicePropositionAny(s, func(e T) bool {
		return e >= v
	})
}

// SlicePropositionNone returns true if the proposition of none of the values in the slice satisfy the proposition.
func SlicePropositionNone[T any](s []T, proposition func(T) bool) bool {
	for _, e := range s {
		if proposition(e) {
			return false
		}
	}
	return true
}

// SlicePropositionAny returns true if the proposition of any of the values in the slice satisfy the proposition.
func SlicePropositionAny[T any](s []T, proposition func(T) bool) bool {
	for _, e := range s {
		if proposition(e) {
			return true
		}
	}
	return false
}

// SlicePropositionAll returns true if the proposition of all of the values in the slice satisfy the proposition.
func SlicePropositionAll[T any](s []T, proposition func(T) bool) bool {
	for _, e := range s {
		if !proposition(e) {
			return false
		}
	}
	return true
}
