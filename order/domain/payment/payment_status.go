package payment

import "github.com/marcosvieirajr/sales-ddd-hexagonal/kernel/errs"

var ErrInvalidPaymentStatus = errs.New("PAYMENT.INVALID_STATUS", "invalid payment status")

// Status represents the lifecycle state of a [Payment].
type Status struct{ value int }

// Define vars for each payment status, starting from 1 to avoid the zero value which can be used as a default or uninitialized state.
var (
	StatusPending    = Status{1} // StatusPending is the initial state; payment is awaiting processing.
	StatusAuthorized = Status{2} // StatusAuthorized indicates the payment was successfully confirmed.
	StatusRefused    = Status{3} // StatusRefused indicates the payment was declined by the gateway.
	StatusRefunded   = Status{4} // StatusRefunded indicates a previously authorized payment was refunded.
	StatusCancelled  = Status{5} // StatusCancelled indicates the payment was cancelled before completion.
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
	return s.value == other.value
}

// ParseStatus converts an int to the corresponding Status value.
// If the input does not match any known status, it returns an error and an empty Status value.
func ParseStatus(value int) (Status, error) {
	s := Status{value}
	if _, ok := statusToString[s]; !ok {
		return Status{}, ErrInvalidPaymentStatus
	}
	return s, nil
}
