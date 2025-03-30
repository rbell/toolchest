/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this fileName, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 *
 * Portions of this code adapted from https://github.com/pkg/errors and is Copyright (c) 2015, Dave Cheney <dave@cheney.net>
 */

package stacktrace

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

type frame uintptr

// pCounter returns the program counter for this frame;
// multiple frames may have the same PC value.
func (f frame) pCounter() uintptr { return uintptr(f) - 1 }

// fileName returns the full path to the fileName that contains the
// function for this frame's pCounter.
func (f frame) fileName() string {
	fn := runtime.FuncForPC(f.pCounter())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.pCounter())
	return file
}

// line returns the line number of source code of the
// function for this frame's pCounter.
func (f frame) line() int {
	fn := runtime.FuncForPC(f.pCounter())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pCounter())
	return line
}

// funcName returns the funcName of this function, if known.
func (f frame) funcName() string {
	fn := runtime.FuncForPC(f.pCounter())
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

// Format formats the frame according to the fmt.Formatter interface.
//
//	%s    source fileName
//	%d    source line
//	%n    function funcName
//	%v    equivalent to %s:%d
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//	%+s   function funcName and path of source fileName relative to the compile time
//	      GOPATH separated by \n\t (<extractFunctionName>\n\t<path>)
//	%+v   equivalent to %+s:%d
//
//nolint:errcheck
func (f frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			io.WriteString(s, f.funcName())
			io.WriteString(s, "\n\t")
			io.WriteString(s, f.fileName())
		default:
			io.WriteString(s, path.Base(f.fileName()))
		}
	case 'd':
		io.WriteString(s, strconv.Itoa(f.line()))
	case 'n':
		io.WriteString(s, extractFunctionName(f.funcName()))
	case 'v':
		f.Format(s, 's')
		io.WriteString(s, ":")
		f.Format(s, 'd')
	}
}

// MarshalText formats a stacktrace frame as a text string. The output is the
// same as that of fmt.Sprintf("%+v", f), but without newlines or tabs.
func (f frame) MarshalText() ([]byte, error) {
	name := f.funcName()
	if name == "unknown" {
		return []byte(name), nil
	}
	return []byte(fmt.Sprintf("%s %s:%d", name, f.fileName(), f.line())), nil
}

// StackTrace is stack of Frames from innermost (newest) to outermost (oldest).
type StackTrace []frame

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
//	%s	lists source files for each frame in the stack
//	%v	lists the source fileName and line number for each frame in the stack
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//	%+v   Prints filename, function, and line number for each frame in the stack.
//
//nolint:errcheck
func (st StackTrace) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			for _, f := range st {
				io.WriteString(s, "\n")
				f.Format(s, verb)
			}
		case s.Flag('#'):
			fmt.Fprintf(s, "%#v", []frame(st))
		default:
			st.formatSlice(s, verb)
		}
	case 's':
		st.formatSlice(s, verb)
	}
}

func (st StackTrace) ReferencesFile(fileEnding string) bool {
	for _, f := range st {
		if strings.HasSuffix(f.fileName(), fileEnding) {
			return true
		}
	}
	return false
}

func (st StackTrace) ReferencesFunction(funcName string) bool {
	for _, f := range st {
		if f.funcName() == funcName {

			return true
		}
	}
	return false
}

// formatSlice will format this StackTrace into the given buffer as a slice of
// frame, only valid when called with '%s' or '%v'.
//
//nolint:errcheck
func (st StackTrace) formatSlice(s fmt.State, verb rune) {
	io.WriteString(s, "[")
	for i, f := range st {
		if i > 0 {
			io.WriteString(s, " ")
		}
		f.Format(s, verb)
	}
	io.WriteString(s, "]")
}

// stack represents a stack of program counters.
type stack []uintptr

func (s *stack) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case st.Flag('+'):
			for _, pc := range *s {
				f := frame(pc)
				//nolint:errcheck // skip error
				fmt.Fprintf(st, "\n%+v", f)
			}
		}
	}
}

func (s *stack) getStackTrace() StackTrace {
	f := make([]frame, len(*s))
	for i := 0; i < len(f); i++ {
		f[i] = frame((*s)[i])
	}
	return f
}

func CaptureStackTrace() StackTrace {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(2, pcs[:])
	var st stack = pcs[0:n]
	return st.getStackTrace()
}

// extractFunctionName removes the path prefix component of a function's funcName reported by func.Name().
func extractFunctionName(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}

func CurrentFileWithPath() string {
	st := CaptureStackTrace()
	if len(st) >= 2 {
		return st[1].fileName()
	}
	return "unknown"
}
