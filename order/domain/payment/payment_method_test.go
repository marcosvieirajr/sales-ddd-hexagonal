package payment_test

import (
	"testing"

	"github.com/marcosvieirajr/sales-ddd-hexagonal/order/domain/payment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMethod_String(t *testing.T) {
	// ==================== Success cases ==================== //
	tests := []struct {
		name   string
		method payment.Method
		want   string
	}{
		{name: "should return 'credit_card' for MethodCreditCard", method: payment.MethodCreditCard, want: "credit_card"},
		{name: "should return 'debit_card' for MethodDebitCard", method: payment.MethodDebitCard, want: "debit_card"},
		{name: "should return 'cash' for MethodCash", method: payment.MethodCash, want: "cash"},
		{name: "should return 'pix' for MethodPix", method: payment.MethodPix, want: "pix"},
		{name: "should return 'bank_transfer' for MethodBankTransfer", method: payment.MethodBankTransfer, want: "bank_transfer"},
		{name: "should return 'banc_slip' for MethodBancSlip", method: payment.MethodBancSlip, want: "banc_slip"},
		// ==================== Failure cases ==================== //
		{name: "should return 'unknown' for zero value (uninitialized)", method: payment.Method{}, want: "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.method.String()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMethod_MarshalText(t *testing.T) {
	tests := []struct {
		name   string
		method payment.Method
		want   string
	}{
		{name: "should marshal MethodCreditCard to 'credit_card'", method: payment.MethodCreditCard, want: "credit_card"},
		{name: "should marshal MethodPix to 'pix'", method: payment.MethodPix, want: "pix"},
		{name: "should marshal unknown method to 'unknown'", method: payment.Method{}, want: "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.method.MarshalText()

			require.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestMethod_Equals(t *testing.T) {
	tests := []struct {
		name   string
		method payment.Method
		other  payment.Method
		want   bool
	}{
		// ==================== Success cases ==================== //
		{name: "should return true when both methods are the same", method: payment.MethodCreditCard, other: payment.MethodCreditCard, want: true},
		{name: "should return true when comparing the same Pix method", method: payment.MethodPix, other: payment.MethodPix, want: true},
		// ==================== Failure cases ==================== //
		{name: "should return false when methods are different", method: payment.MethodCreditCard, other: payment.MethodDebitCard, want: false},
		{name: "should return false when comparing with an uninitialized method", method: payment.MethodCreditCard, other: payment.Method{}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.method.Equals(tt.other)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseMethod(t *testing.T) {
	// ==================== Success cases ==================== //
	successTests := []struct {
		name       string
		value      int
		wantMethod payment.Method
	}{
		{name: "should parse 1 to MethodCreditCard", value: 1, wantMethod: payment.MethodCreditCard},
		{name: "should parse 2 to MethodDebitCard", value: 2, wantMethod: payment.MethodDebitCard},
		{name: "should parse 3 to MethodCash", value: 3, wantMethod: payment.MethodCash},
		{name: "should parse 4 to MethodPix", value: 4, wantMethod: payment.MethodPix},
		{name: "should parse 5 to MethodBankTransfer", value: 5, wantMethod: payment.MethodBankTransfer},
		{name: "should parse 6 to MethodBancSlip", value: 6, wantMethod: payment.MethodBancSlip},
	}
	for _, tt := range successTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := payment.ParseMethod(tt.value)

			require.NoError(t, err)
			assert.Equal(t, tt.wantMethod, got)
		})
	}

	// ==================== Failure cases ==================== //
	failureTests := []struct {
		name    string
		value   int
		wantErr error
	}{
		{name: "should return an error for zero value (uninitialized)", value: 0, wantErr: payment.ErrInvalidPaymentMethod},
		{name: "should return an error for a negative value", value: -1, wantErr: payment.ErrInvalidPaymentMethod},
		{name: "should return an error for an out-of-range value", value: 999, wantErr: payment.ErrInvalidPaymentMethod},
	}
	for _, tt := range failureTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := payment.ParseMethod(tt.value)

			require.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, payment.Method{}, got)
		})
	}
}
