package domain

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

// GenerateID returns a new unique identifier string for domain entities.
// TODO: replace the stub implementation with a proper generator (e.g. UUID or ULID).
func GenerateID() string {
	return "12345678-1234-5678-1234-567812345678"
}
