package order

import "github.com/marcosvieirajr/sales-ddd-hexagonal/kernel/errs"

var ErrInvalidOrderStatus = errs.New("ORDER.INVALID_STATUS", "invalid order status")

// Status represents the fulfillment lifecycle state of an [Order].
type Status struct{ value int }

var (
	StatusCreated    = Status{1} // StatusCreated is the initial state of an order after placement.
	StatusPaid       = Status{2} // StatusPaid indicates the order payment has been confirmed.
	StatusSeparating = Status{3} // StatusSeparating indicates the order is being picked and packed.
	StatusShipped    = Status{4} // StatusShipped indicates the order has been dispatched to the carrier.
	StatusDelivered  = Status{5} // StatusDelivered indicates the order has reached the customer.
	StatusCancelled  = Status{6} // StatusCancelled indicates the order has been cancelled.
)

var statusToString = map[Status]string{
	StatusCreated:    "created",
	StatusPaid:       "paid",
	StatusSeparating: "separating",
	StatusShipped:    "shipped",
	StatusDelivered:  "delivered",
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
		return Status{}, ErrInvalidOrderStatus
	}
	return s, nil
}
