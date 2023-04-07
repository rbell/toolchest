/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package errors

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// ValidationError is error suitable for reporting validation errors to a user.  Errors reported can be warnings or errors.  Supports validation against nested types (i.e. Customer has an address).
type ValidationError struct {
	errorMap   map[string][]string
	warningMap map[string][]string
	// children contains errors for fields that reference another validated structure
	children map[string]*ValidationError
}

// NewValidationError returns a new validation error
func NewValidationError(context, msg string, isWarning bool) *ValidationError {
	em := make(map[string][]string)
	em[context] = []string{msg}
	if isWarning {
		return &ValidationError{warningMap: em, errorMap: make(map[string][]string)}
	}
	return &ValidationError{errorMap: em, warningMap: make(map[string][]string)}
}

// NewValidationErrors returns a new validation error for a map of error messages
func NewValidationErrors(errs map[string][]string, children map[string]*ValidationError) *ValidationError {
	if errs == nil {
		return &ValidationError{errorMap: make(map[string][]string), children: children}
	}
	return &ValidationError{errorMap: errs, children: children}
}

// NewValidationErrorsWithWarnings returns a new validation error for a map of error messages
func NewValidationErrorsWithWarnings(errs, warnings map[string][]string, children map[string]*ValidationError) *ValidationError {
	if errs == nil && warnings == nil {
		return &ValidationError{errorMap: make(map[string][]string), children: children}
	}
	return &ValidationError{errorMap: errs, warningMap: warnings, children: children}
}

// IsValidatorError returns reference to ValidationError and a bool indicating if the err passed in is a ValidationError
func IsValidatorError(err error) (*ValidationError, bool) {
	if e, ok := err.(*ValidationError); ok {
		return e, true
	}
	return nil, false
}

// AddErrorToValidation joins two errors together into a ValidatorError
func AddErrorToValidation(e1, e2 error) *ValidationError {
	var e *ValidationError
	// if e1 is nil and e2 is not, short circuit using e2
	if (e1 == nil || reflect.ValueOf(e1).IsNil()) && e2 != nil {
		if errors.As(e2, &e) {
			// e1 is nil and e2 is a ValidationError, return e2
			return e2.(*ValidationError)
		}
		// e1 is nil and e2 is an error.  Return new validation error using
		return NewValidationError("", e2.Error(), false)
	}

	var ve *ValidationError
	if errors.As(e1, &e) {
		//nolint:errcheck // above line infers its castable
		ve = e1.(*ValidationError)
	} else {
		ve = NewValidationError("", e1.Error(), false)
	}

	if errors.As(e2, &e) {
		errMap := ve.GetErrorMap()
		for key, msg := range e2.(*ValidationError).GetFlatErrorMap() {
			addMsgs(errMap, key, msg...)
		}
		warnMap := ve.GetWarningMap()
		for key, msg := range e2.(*ValidationError).GetFlatWarningMap() {
			addMsgs(warnMap, key, msg...)
		}
		childErrs := ve.GetChildErrors()
		for key, ve := range e2.(*ValidationError).GetChildErrors() {
			childErrs[key] = ve
		}
	} else {
		errMap := ve.GetErrorMap()
		addMsgs(errMap, "", e2.Error())
	}

	return ve
}

// Error returns the error messages in a single string
func (e *ValidationError) Error() string {
	sb := strings.Builder{}
	for _, ee := range e.GetFlatErrorMap() {
		for _, e := range ee {
			sb.WriteString(fmt.Sprintf("ERROR: %v\n", e))
		}
	}
	for _, ee := range e.GetFlatWarningMap() {
		for _, e := range ee {
			sb.WriteString(fmt.Sprintf("WARNING: %v\n", e))
		}
	}
	return sb.String()
}

// GetFlatErrorMap gets a map of error messages mapped by field
func (e *ValidationError) GetFlatErrorMap() map[string][]string {
	flatMap := e.errorMap
	for k, v := range e.children {
		childErrors := getFlattenedMap(k, v, false)
		for ek, e := range childErrors {
			addMsgs(flatMap, ek, e...)
		}
	}
	return flatMap
}

// GetFlatWarningMap gets a map of warning messages mapped by field
func (e *ValidationError) GetFlatWarningMap() map[string][]string {
	flatMap := e.warningMap
	for k, v := range e.children {
		childErrors := getFlattenedMap(k, v, true)
		for ek, e := range childErrors {
			addMsgs(flatMap, ek, e...)
		}
	}
	return flatMap
}

// GetErrorMap returns the top level error messages
func (e *ValidationError) GetErrorMap() map[string][]string {
	return e.errorMap
}

// GetWarningMap returns the top level warning messages
func (e *ValidationError) GetWarningMap() map[string][]string {
	return e.warningMap
}

// GetChildErrors returns errors of structs referenced by the struct being validated
func (e *ValidationError) GetChildErrors() map[string]*ValidationError {
	return e.children
}

func getFlattenedMap(key string, ve *ValidationError, getWarnings bool) map[string][]string {
	flatMap := prefixKeys(ve.errorMap, key+".")
	if getWarnings {
		flatMap = prefixKeys(ve.warningMap, key+".")
	}
	for childKey, child := range ve.children {
		// resursively call geFlattenedMap
		childMap := getFlattenedMap(key+"."+childKey, child, getWarnings)
		for k, e := range childMap {
			addMsgs(flatMap, k, e...)
		}
	}
	return flatMap
}

func prefixKeys(m map[string][]string, prefix string) map[string][]string {
	result := make(map[string][]string)
	for k, v := range m {
		result[prefix+k] = v
	}
	return result
}

func addMsgs(errMap map[string][]string, context string, msg ...string) {
	if _, ok := errMap[context]; !ok {
		errMap[context] = []string{}
	}
	errMap[context] = append(errMap[context], msg...)
}
