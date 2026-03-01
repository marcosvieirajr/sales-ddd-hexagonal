package payment_test

import (
	"testing"

	"github.com/marcosvieirajr/sales-ddd-hexagonal/order/domain/payment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatus_String(t *testing.T) {
	// ==================== Success cases ==================== //
	tests := []struct {
		name   string
		status payment.Status
		want   string
	}{
		{name: "should return 'pending' for StatusPending", status: payment.StatusPending, want: "pending"},
		{name: "should return 'authorized' for StatusAuthorized", status: payment.StatusAuthorized, want: "authorized"},
		{name: "should return 'refused' for StatusRefused", status: payment.StatusRefused, want: "refused"},
		{name: "should return 'refunded' for StatusRefunded", status: payment.StatusRefunded, want: "refunded"},
		{name: "should return 'cancelled' for StatusCancelled", status: payment.StatusCancelled, want: "cancelled"},
		// ==================== Failure cases ==================== //
		{name: "should return 'unknown' for an unrecognized status value", status: payment.Status{}, want: "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.String()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStatus_MarshalText(t *testing.T) {
	tests := []struct {
		name   string
		status payment.Status
		want   string
	}{
		{name: "should marshal StatusPending to 'pending'", status: payment.StatusPending, want: "pending"},
		{name: "should marshal StatusAuthorized to 'authorized'", status: payment.StatusAuthorized, want: "authorized"},
		{name: "should marshal unknown status to 'unknown'", status: payment.Status{}, want: "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.status.MarshalText()

			require.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestStatus_Equals(t *testing.T) {
	tests := []struct {
		name   string
		status payment.Status
		other  payment.Status
		want   bool
	}{
		// ==================== Success cases ==================== //
		{name: "should return true when both statuses are the same", status: payment.StatusPending, other: payment.StatusPending, want: true},
		{name: "should return true when comparing the same Authorized status", status: payment.StatusAuthorized, other: payment.StatusAuthorized, want: true},
		// ==================== Failure cases ==================== //
		{name: "should return false when statuses are different", status: payment.StatusPending, other: payment.StatusAuthorized, want: false},
		{name: "should return false when comparing with an uninitialized status", status: payment.StatusPending, other: payment.Status{}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.Equals(tt.other)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseStatus(t *testing.T) {
	// ==================== Success cases ==================== //
	successTests := []struct {
		name       string
		value      int
		wantStatus payment.Status
	}{
		{name: "should parse 1 to StatusPending", value: 1, wantStatus: payment.StatusPending},
		{name: "should parse 2 to StatusAuthorized", value: 2, wantStatus: payment.StatusAuthorized},
		{name: "should parse 3 to StatusRefused", value: 3, wantStatus: payment.StatusRefused},
		{name: "should parse 4 to StatusRefunded", value: 4, wantStatus: payment.StatusRefunded},
		{name: "should parse 5 to StatusCancelled", value: 5, wantStatus: payment.StatusCancelled},
	}
	for _, tt := range successTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := payment.ParseStatus(tt.value)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, got)
		})
	}

	// ==================== Failure cases ==================== //
	failureTests := []struct {
		name    string
		value   int
		wantErr error
	}{
		{name: "should return an error for a negative value", value: -1, wantErr: payment.ErrInvalidPaymentStatus},
		{name: "should return an error for an out-of-range value", value: 999, wantErr: payment.ErrInvalidPaymentStatus},
	}
	for _, tt := range failureTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := payment.ParseStatus(tt.value)

			require.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, payment.Status{}, got)
		})
	}
}
