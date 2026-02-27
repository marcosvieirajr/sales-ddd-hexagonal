package domain

import (
	"reflect"
	"regexp"
	"strings"
)

// CheckMatchRegex returns err if value does not match the regular expression regex,
// or nil when the value matches.
func CheckMatchRegex(value string, regex *regexp.Regexp, err error) error {
	if !regex.MatchString(value) {
		return err
	}
	return nil
}

// CheckNotNullOrWhiteSpace returns err if value is empty or contains only whitespace,
// or nil when value contains at least one non-whitespace character.
func CheckNotNullOrWhiteSpace(value string, err error) error {
	if strings.TrimSpace(value) == "" {
		return err
	}
	return nil
}

// CheckNotZeroOrNegative returns err if value is zero or negative (≤ 0),
// or nil when value is strictly positive.
func CheckNotZeroOrNegative(value float64, err error) error {
	if value <= 0 {
		return err
	}
	return nil
}

// CheckNotNil returns err if value is nil, or nil when value is non-nil.
// It is the inverse of [CheckNil] and is intended for validating pointer or interface
// fields that must be set (e.g. a required transaction code).
// It handles typed nil pointers correctly by inspecting the underlying reflect value.
func CheckNotNil(value any, err error) error {
	if isNil(value) {
		return err
	}
	return nil
}

// CheckNil returns err if value is non-nil, or nil when value is nil.
// It is the inverse of [CheckNotNil] and is intended for validating pointer or interface
// fields that must not already be set (e.g. preventing a second assignment).
// It handles typed nil pointers correctly by inspecting the underlying reflect value.
func CheckNil(value any, err error) error {
	if !isNil(value) {
		return err
	}
	return nil
}

// isNil reports whether value is nil, handling both untyped nil interfaces and
// typed nil pointers (e.g. (*string)(nil) passed as any).
func isNil(value any) bool {
	if value == nil {
		return true
	}
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		return rv.IsNil()
	default:
		return false
	}
}
