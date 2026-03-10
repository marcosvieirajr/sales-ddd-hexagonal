package order_test

import (
	"testing"

	order "github.com/marcosvieirajr/sales-ddd-hexagonal/order/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatus_String(t *testing.T) {
	// ==================== Success cases ==================== //
	tests := []struct {
		name   string
		status order.Status
		want   string
	}{
		{name: "should return 'pending' for StatusPending", status: order.StatusPending, want: "pending"},
		{name: "should return 'paid' for StatusPaid", status: order.StatusPaid, want: "paid"},
		{name: "should return 'separating' for StatusSeparating", status: order.StatusSeparating, want: "separating"},
		{name: "should return 'shipped' for StatusShipped", status: order.StatusShipped, want: "shipped"},
		{name: "should return 'delivered' for StatusDelivered", status: order.StatusDelivered, want: "delivered"},
		{name: "should return 'cancelled' for StatusCancelled", status: order.StatusCancelled, want: "cancelled"},
		// ==================== Failure cases ==================== //
		{name: "should return 'unknown' for an unrecognized status value", status: order.Status{}, want: "unknown"},
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
		status order.Status
		want   string
	}{
		// ==================== Success cases ==================== //
		{name: "should marshal StatusPending to 'pending'", status: order.StatusPending, want: "pending"},
		{name: "should marshal StatusPaid to 'paid'", status: order.StatusPaid, want: "paid"},
		// ==================== Failure cases ==================== //
		{name: "should marshal unknown status to 'unknown'", status: order.Status{}, want: "unknown"},
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
		status order.Status
		other  order.Status
		want   bool
	}{
		// ==================== Success cases ==================== //
		{name: "should return true when both statuses are the same", status: order.StatusPending, other: order.StatusPending, want: true},
		{name: "should return true when comparing the same Paid status", status: order.StatusPaid, other: order.StatusPaid, want: true},
		// ==================== Failure cases ==================== //
		{name: "should return false when statuses are different", status: order.StatusPending, other: order.StatusPaid, want: false},
		{name: "should return false when comparing with an uninitialized status", status: order.StatusPending, other: order.Status{}, want: false},
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
		wantStatus order.Status
	}{
		{name: "should parse 1 to StatusPending", value: 1, wantStatus: order.StatusPending},
		{name: "should parse 2 to StatusPaid", value: 2, wantStatus: order.StatusPaid},
		{name: "should parse 3 to StatusSeparating", value: 3, wantStatus: order.StatusSeparating},
		{name: "should parse 4 to StatusShipped", value: 4, wantStatus: order.StatusShipped},
		{name: "should parse 5 to StatusDelivered", value: 5, wantStatus: order.StatusDelivered},
		{name: "should parse 6 to StatusCancelled", value: 6, wantStatus: order.StatusCancelled},
	}
	for _, tt := range successTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := order.ParseStatus(tt.value)

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
		{name: "should return an error for a negative value", value: -1, wantErr: order.ErrInvalidOrderStatus},
		{name: "should return an error for an out-of-range value", value: 999, wantErr: order.ErrInvalidOrderStatus},
	}
	for _, tt := range failureTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := order.ParseStatus(tt.value)

			require.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, order.Status{}, got)
		})
	}
}
