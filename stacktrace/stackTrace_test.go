/*
 * Copyright (c) 2025 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package stacktrace

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrentFile_ReturnsStringContainingThisFile(t *testing.T) {
	fileName := CurrentFileWithPath()
	assert.Contains(t, fileName, "stackTrace_test.go")
}

func TestStackTrace_ReferenceFile_ReturnsTrue(t *testing.T) {
	st := CaptureStackTrace()
	references := st.ReferencesFile(`testing/testing.go`)
	assert.True(t, references)
}

func TestStackTrace_ReferencesFunction(t *testing.T) {
	st := CaptureStackTrace()
	references := st.ReferencesFunction("testing.tRunner")
	assert.True(t, references)
}

func TestStackTrace_Format_WithFlags(t *testing.T) {
	st := CaptureStackTrace()
	result := fmt.Sprintf("%+v", st)
	expected := `
github.com/rbell/toolchest/stacktrace.TestStackTrace_Format_WithFlags
	/home/rbell/dev/github.com/rbell/toolchest/stacktrace/stackTrace_test.go:34
testing.tRunner
	/home/rbell/sdk/go1.22.3/src/testing/testing.go:1689
runtime.goexit
	/home/rbell/sdk/go1.22.3/src/runtime/asm_amd64.s:1695`
	assert.Equal(t, expected, result)
}

func TestStackTrace_Format_Default(t *testing.T) {
	st := CaptureStackTrace()
	result := fmt.Sprint(st)
	expected := `[stackTrace_test.go:47 testing.go:1689 asm_amd64.s:1695]`
	assert.Equal(t, expected, result)
}
