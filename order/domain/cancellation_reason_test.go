package order_test

import (
	"testing"

	order "github.com/marcosvieirajr/sales-ddd-hexagonal/order/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCancellationReason_String(t *testing.T) {
	// ==================== Success cases ==================== //
	tests := []struct {
		name   string
		reason order.CancellationReason
		want   string
	}{
		{name: "should return 'customer_cancelled' for CancellationReasonCustomerCancelled", reason: order.CancellationReasonCustomerCancelled, want: "customer_cancelled"},
		{name: "should return 'payment_error' for CancellationReasonPaymentError", reason: order.CancellationReasonPaymentError, want: "payment_error"},
		{name: "should return 'out_of_stock' for CancellationReasonOutOfStock", reason: order.CancellationReasonOutOfStock, want: "out_of_stock"},
		{name: "should return 'invalid_address' for CancellationReasonInvalidAddress", reason: order.CancellationReasonInvalidAddress, want: "invalid_address"},
		{name: "should return 'other' for CancellationReasonOther", reason: order.CancellationReasonOther, want: "other"},
		// ==================== Failure cases ==================== //
		{name: "should return 'unknown' for an unrecognized reason value", reason: order.CancellationReason{}, want: "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.reason.String()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCancellationReason_MarshalText(t *testing.T) {
	tests := []struct {
		name   string
		reason order.CancellationReason
		want   string
	}{
		// ==================== Success cases ==================== //
		{name: "should marshal CancellationReasonCustomerCancelled to 'customer_cancelled'", reason: order.CancellationReasonCustomerCancelled, want: "customer_cancelled"},
		{name: "should marshal CancellationReasonPaymentError to 'payment_error'", reason: order.CancellationReasonPaymentError, want: "payment_error"},
		// ==================== Failure cases ==================== //
		{name: "should marshal unknown reason to 'unknown'", reason: order.CancellationReason{}, want: "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.reason.MarshalText()

			require.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestCancellationReason_Equals(t *testing.T) {
	tests := []struct {
		name   string
		reason order.CancellationReason
		other  order.CancellationReason
		want   bool
	}{
		// ==================== Success cases ==================== //
		{name: "should return true when both reasons are the same", reason: order.CancellationReasonCustomerCancelled, other: order.CancellationReasonCustomerCancelled, want: true},
		{name: "should return true when comparing the same PaymentError reason", reason: order.CancellationReasonPaymentError, other: order.CancellationReasonPaymentError, want: true},
		// ==================== Failure cases ==================== //
		{name: "should return false when reasons are different", reason: order.CancellationReasonCustomerCancelled, other: order.CancellationReasonPaymentError, want: false},
		{name: "should return false when comparing with an uninitialized reason", reason: order.CancellationReasonCustomerCancelled, other: order.CancellationReason{}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.reason.Equals(tt.other)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseCancellationReason(t *testing.T) {
	// ==================== Success cases ==================== //
	successTests := []struct {
		name       string
		value      int
		wantReason order.CancellationReason
	}{
		{name: "should parse 1 to CancellationReasonCustomerCancelled", value: 1, wantReason: order.CancellationReasonCustomerCancelled},
		{name: "should parse 2 to CancellationReasonPaymentError", value: 2, wantReason: order.CancellationReasonPaymentError},
		{name: "should parse 3 to CancellationReasonOutOfStock", value: 3, wantReason: order.CancellationReasonOutOfStock},
		{name: "should parse 4 to CancellationReasonInvalidAddress", value: 4, wantReason: order.CancellationReasonInvalidAddress},
		{name: "should parse 5 to CancellationReasonOther", value: 5, wantReason: order.CancellationReasonOther},
	}
	for _, tt := range successTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := order.ParseCancellationReason(tt.value)

			require.NoError(t, err)
			assert.Equal(t, tt.wantReason, got)
		})
	}

	// ==================== Failure cases ==================== //
	failureTests := []struct {
		name    string
		value   int
		wantErr error
	}{
		{name: "should return an error for a negative value", value: -1, wantErr: order.ErrInvalidCancellationReason},
		{name: "should return an error for an out-of-range value", value: 999, wantErr: order.ErrInvalidCancellationReason},
	}
	for _, tt := range failureTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := order.ParseCancellationReason(tt.value)

			require.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, order.CancellationReason{}, got)
		})
	}
}
