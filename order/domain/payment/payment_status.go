package payment

import "github.com/marcosvieirajr/sales-ddd-hexagonal/kernel/errs"

var ErrInvalidPaymentStatus = errs.New("PAYMENT.INVALID_STATUS", "invalid payment status")

// Status represents the lifecycle state of a [Payment].
type Status int

// Define constants for each payment status, starting from 0 to use the zero value as the initial pending state.
const (
	StatusPending    Status = iota + 1 // StatusPending is the initial state; payment is awaiting processing.
	StatusAuthorized                   // StatusAuthorized indicates the payment was successfully confirmed.
	StatusRefused                      // StatusRefused indicates the payment was declined by the gateway.
	StatusRefunded                     // StatusRefunded indicates a previously authorized payment was refunded.
	StatusCancelled                    // StatusCancelled indicates the payment was cancelled before completion.
)

// statusToString maps Status values to their string representations.
var statusToString = map[Status]string{
	StatusPending:    "pending",
	StatusAuthorized: "authorized",
	StatusRefused:    "refused",
	StatusRefunded:   "refunded",
	StatusCancelled:  "cancelled",
}

// String returns the string representation of the Status.
func (s Status) String() string {
	if str, ok := statusToString[s]; ok {
		return str
	}
	return "unknown"
}

// MarshalText provides support for logging and any marshal needs.
func (s Status) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// Equals checks if two Status values are equal.
func (s Status) Equals(other Status) bool {
	return s == other
}

// ParseStatus converts an int to the corresponding Status value.
// If the input does not match any known status, it returns an error and an empty Status value.
func ParseStatus(value int) (Status, error) {
	status := Status(value)
	if _, ok := statusToString[status]; !ok {
		return 0, ErrInvalidPaymentStatus
	}
	return status, nil
}
