package order

import "github.com/marcosvieirajr/sales-ddd-hexagonal/kernel/errs"

var ErrInvalidCancellationReason = errs.New("ORDER.INVALID_CANCELLATION_REASON", "invalid cancellation reason")

// CancellationReason represents the fulfillment cancellation reason or and [Order].
type CancellationReason struct {
	value int
}

var (
	CancellationReasonCustomerCancelled = CancellationReason{1}
	CancellationReasonPaymentError      = CancellationReason{2}
	CancellationReasonOutOfStock        = CancellationReason{3}
	CancellationReasonInvalidAddress    = CancellationReason{4}
	CancellationReasonOther             = CancellationReason{5}
)

var cancellationToString = map[CancellationReason]string{
	CancellationReasonCustomerCancelled: "customer_cancelled",
	CancellationReasonPaymentError:      "payment_error",
	CancellationReasonOutOfStock:        "out_of_stock",
	CancellationReasonInvalidAddress:    "invalid_address",
	CancellationReasonOther:             "other",
}

// String returns the string representation of the CancellationReason.
func (s CancellationReason) String() string {
	if str, ok := cancellationToString[s]; ok {
		return str
	}
	return "unknown"
}

// MarshalText provides support for logging and any marshal needs.
func (s CancellationReason) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// Equals checks if two CancellationReason values are equal.
func (s CancellationReason) Equals(other CancellationReason) bool {
	return s.value == other.value
}

// ParseCancellationReason converts an int to the corresponding CancellationReason value.
// If the input does not match any known cancellation reason, it returns an error and an empty CancellationReason value.
func ParseCancellationReason(value int) (CancellationReason, error) {
	r := CancellationReason{value: value}
	if _, ok := cancellationToString[r]; !ok {
		return CancellationReason{}, ErrInvalidCancellationReason
	}
	return r, nil
}
