package payment

// Status represents the lifecycle state of a [Payment].
type Status int

const (
	StatusPending    Status = iota // StatusPending is the initial state; payment is awaiting processing.
	StatusAuthorized               // StatusAuthorized indicates the payment was successfully confirmed.
	StatusRefused                  // StatusRefused indicates the payment was declined by the gateway.
	StatusRefunded                 // StatusRefunded indicates a previously authorized payment was refunded.
	StatusCancelled                // StatusCancelled indicates the payment was cancelled before completion.
)

// Equals reports whether s and other represent the same payment status.
func (s Status) Equals(other Status) bool {
	return s == other
}
