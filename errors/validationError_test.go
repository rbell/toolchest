/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package errors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidationError_GetErrorMap_NoHierarchy_ReturnsFlatMap(t *testing.T) {
	// Setup
	ve := &ValidationError{
		errorMap: map[string][]string{
			"TestField": {
				"Err1",
				"Err2",
			},
		},
	}

	// Test
	result := ve.GetFlatErrorMap()

	// Assert
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Contains(t, result, "TestField")
	assert.Len(t, result["TestField"], 2)
}

func TestValidationError_GetErrorMap_OneChild_ReturnsFlatMap(t *testing.T) {
	// Setup
	ve := &ValidationError{
		errorMap: map[string][]string{
			"TestField": {
				"Err1",
				"Err2",
			},
		},
		children: map[string]*ValidationError{
			"TestRef": {
				errorMap: map[string][]string{
					"ChildTestField": {
						"ChildErr1",
						"ChildErr2",
						"ChildErr3",
					},
				},
			},
		},
	}

	// Test
	result := ve.GetFlatErrorMap()

	// Assert
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Contains(t, result, "TestField")
	assert.Len(t, result["TestField"], 2)
	assert.Contains(t, result, "TestRef.ChildTestField")
	assert.Len(t, result["TestRef.ChildTestField"], 3)
}

func TestValidationError_GetErrorMap_OneGrandChild_ReturnsFlatMap(t *testing.T) {
	// Setup
	ve := &ValidationError{
		errorMap: map[string][]string{
			"TestField": {
				"Err1",
				"Err2",
			},
		},
		children: map[string]*ValidationError{
			"TestRef": {
				errorMap: map[string][]string{
					"ChildTestField": {
						"ChildErr1",
						"ChildErr2",
						"ChildErr3",
					},
				},
				children: map[string]*ValidationError{
					"GrandChildTestRef": {
						errorMap: map[string][]string{
							"GrandChildTestField": {
								"GrandChildErr1",
								"GrandChildErr2",
								"GrandChildErr3",
								"GrandChildErr4",
							},
						},
					},
				},
			},
		},
	}

	// Test
	result := ve.GetFlatErrorMap()

	// Assert
	assert.NotNil(t, result)
	assert.Len(t, result, 3)
	assert.Contains(t, result, "TestField")
	assert.Len(t, result["TestField"], 2)
	assert.Contains(t, result, "TestRef.ChildTestField")
	assert.Len(t, result["TestRef.ChildTestField"], 3)
	assert.Contains(t, result, "TestRef.GrandChildTestRef.GrandChildTestField")
	assert.Len(t, result["TestRef.GrandChildTestRef.GrandChildTestField"], 4)
}
