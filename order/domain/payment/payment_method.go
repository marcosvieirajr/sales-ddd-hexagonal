package payment

import "github.com/marcosvieirajr/sales-ddd-hexagonal/kernel/errs"

var ErrInvalidPaymentMethod = errs.New("PAYMENT.INVALID_METHOD", "invalid payment method")

// Method represents the payment method chosen by the customer.
type Method struct{ value int }

// Define vars for each payment method, starting from 1 to avoid the zero value which can be used as a default or uninitialized state.
var (
	MethodCreditCard   = Method{1} // MethodCreditCard represents payment by credit card.
	MethodDebitCard    = Method{2} // MethodDebitCard represents payment by debit card.
	MethodCash         = Method{3} // MethodCash represents payment in cash.
	MethodPix          = Method{4} // MethodPix represents payment via Pix instant transfer.
	MethodBankTransfer = Method{5} // MethodBankTransfer represents payment via bank transfer (TED/DOC).
	MethodBancSlip     = Method{6} // MethodBancSlip represents payment via bank slip (boleto bancário).
)

// methodToString maps Method values to their string representations.
var methodToString = map[Method]string{
	MethodCreditCard:   "credit_card",
	MethodDebitCard:    "debit_card",
	MethodCash:         "cash",
	MethodPix:          "pix",
	MethodBankTransfer: "bank_transfer",
	MethodBancSlip:     "banc_slip",
}

// String returns the string representation of the Method.
func (m Method) String() string {
	if str, ok := methodToString[m]; ok {
		return str
	}
	return "unknown"
}

// MarshalText provides support for logging and any marshal needs.
func (m Method) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

// Equals checks if two Method values are equal.
func (m Method) Equals(other Method) bool {
	return m.value == other.value
}

// ParseMethod converts an int to the corresponding Method value.
// If the input does not match any known method, it returns an error and an empty Method value.
func ParseMethod(value int) (Method, error) {
	m := Method{value}
	if _, ok := methodToString[m]; !ok {
		return Method{}, ErrInvalidPaymentMethod
	}
	return m, nil
}
