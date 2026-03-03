package kernel

import (
	"github.com/oklog/ulid/v2"
)

// Must is a convenience generic function that returns x if err is nil,
// or panics with err otherwise. It is intended for use in package-level
// variable initialization where a non-nil error indicates an unrecoverable
// programming mistake rather than an expected runtime failure.
func Must[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}
	return x
}

// NewID returns a new unique, lexicographically sortable, and
// monotonically increasing identifier. Safe for concurrent use.
func NewID() ulid.ULID {
	return ulid.Make()
}
