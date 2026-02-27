// Package errs provides domain-specific error types for the sales domain.
// It defines [DomainError], a structured error that carries an [ErrorCode] and a
// human-readable message, supporting sentinel-based error matching via [errors.Is].
package errs

import (
	"errors"
	"fmt"
)

// ErrorCode is a string identifier for a domain error.
// By convention it follows the SCREAMING_SNAKE_CASE pattern "AGGREGATE.REASON"
// (e.g. "ORDER_ITEM.NEGATIVE_DISCOUNT").
type ErrorCode string

// DomainError represents a business rule or domain invariant violation.
// It carries a structured [ErrorCode] for programmatic matching and a human-readable
// Message for logging or display. An optional Err field allows wrapping lower-level
// errors into the domain error chain.
type DomainError struct {
	Code    ErrorCode // e.g. "ORDER_ITEM.NEGATIVE_DISCOUNT"
	Message string    // human-readable description of the violation
	Err     error     // optional underlying error for wrapping
}

// Error returns a formatted string representation of the error.
// When an underlying error is present it is appended after a colon:
// "[CODE] message: underlying error". Otherwise it returns "[CODE] message".
func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error so that [errors.Is] and [errors.As]
// can traverse the full error chain.
func (e *DomainError) Unwrap() error {
	return e.Err
}

// Is reports whether e has the same [ErrorCode] as target.
// This enables sentinel [DomainError] values to match via [errors.Is] even
// when the returned error has been wrapped or copied, as comparison is done
// by Code rather than by pointer identity.
func (e *DomainError) Is(target error) bool {
	var domErr *DomainError
	ok := errors.As(target, &domErr)
	if !ok {
		return false
	}
	return e.Code == domErr.Code
}

// Wrap returns a shallow copy of e with Err set to err.
// The copy preserves the original Code and Message, while [errors.Unwrap]
// will traverse to err. Use this to attach a lower-level cause to a sentinel error.
func (e *DomainError) Wrap(err error) *DomainError {
	return &DomainError{Code: e.Code, Message: e.Message, Err: err}
}

// New creates a [DomainError] with the given code and human-readable message.
// Use this to define package-level sentinel errors for domain invariant violations.
func New(code ErrorCode, message string) *DomainError {
	return &DomainError{Code: code, Message: message}
}

// Wrap creates a [DomainError] with the given code and message, wrapping err
// as the underlying cause. Use this when a domain rule violation originates
// from a lower-level error that should remain accessible via [errors.Unwrap].
func Wrap(code ErrorCode, message string, err error) *DomainError {
	return &DomainError{Code: code, Message: message, Err: err}
}
