package payment_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/marcosvieirajr/sales-ddd-hexagonal/order/domain/payment"
	"github.com/marcosvieirajr/sales-ddd-hexagonal/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createValidPayment(t *testing.T) *payment.Payment {
	t.Helper()
	return shared.Must(payment.NewPayment("order-123", 100.0, payment.MethodCreditCard))
}

func createPaymentWithCode(t *testing.T) *payment.Payment {
	t.Helper()
	p := createValidPayment(t)
	require.NoError(t, p.DefineTransactionCode("TXN-123"))
	return p
}

func TestNewPayment(t *testing.T) {
	t.Run("should successfully create a new payment with valid input", func(t *testing.T) {
		got, err := payment.NewPayment("order-123", 100.0, payment.MethodCreditCard)

		require.NoError(t, err)
		want := &payment.Payment{
			OrderID: "order-123",
			Amount:  100.0,
			Method:  payment.MethodCreditCard,
			Status:  payment.StatusPending,
		}
		ignoreFields := cmpopts.IgnoreFields(payment.Payment{}, "ID") // ignore ID since it's generated and not predictable
		assert.True(t, cmp.Equal(got, want, ignoreFields), "got and want should be equal ignoring ID: %v", cmp.Diff(got, want, ignoreFields))
	})

	t.Run("should return an error when invalid input is provided", func(t *testing.T) {
		type args struct {
			orderID string
			amount  float64
			method  payment.Method
		}
		tests := []struct {
			name    string
			args    args
			wantErr error
		}{
			{
				name:    "should return an error when order ID is empty",
				args:    args{orderID: "", amount: 100.0, method: payment.MethodCreditCard},
				wantErr: payment.ErrInvalidOrderID,
			},
			{
				name:    "should return an error when order ID is whitespace",
				args:    args{orderID: "   ", amount: 100.0, method: payment.MethodCreditCard},
				wantErr: payment.ErrInvalidOrderID,
			},
			{
				name:    "should return an error when amount is zero",
				args:    args{orderID: "order-123", amount: 0.0, method: payment.MethodCreditCard},
				wantErr: payment.ErrInvalidPaymentAmount,
			},
			{
				name:    "should return an error when amount is negative",
				args:    args{orderID: "order-123", amount: -10.0, method: payment.MethodCreditCard},
				wantErr: payment.ErrInvalidPaymentAmount,
			},
			{
				name:    "should return an error for invalid order ID when both fields are invalid",
				args:    args{orderID: "", amount: 0.0, method: payment.MethodCreditCard},
				wantErr: payment.ErrInvalidOrderID,
			},
			{
				name:    "should return an error for invalid amount when both fields are invalid",
				args:    args{orderID: "", amount: 0.0, method: payment.MethodCreditCard},
				wantErr: payment.ErrInvalidPaymentAmount,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := payment.NewPayment(tt.args.orderID, tt.args.amount, tt.args.method)

				require.Nil(t, got)
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			})
		}
	})
}

func TestPayment_DefineTransactionCode(t *testing.T) {
	t.Run("should successfully define transaction code with valid code", func(t *testing.T) {
		p := createValidPayment(t)

		err := p.DefineTransactionCode("TXN-123")

		require.NoError(t, err)
		assert.Equal(t, "TXN-123", *p.TransactionCode)
		assert.NotNil(t, p.UpdatedAt, "UpdatedAt should be set on success")
	})

	t.Run("should return an error when input is invalid", func(t *testing.T) {
		tests := []struct {
			name    string
			setup   func(t *testing.T) *payment.Payment
			code    string
			wantErr error
		}{
			{
				name:    "should return an error when code is empty",
				setup:   func(t *testing.T) *payment.Payment { return createValidPayment(t) },
				code:    "",
				wantErr: payment.ErrInvalidTransactionCode,
			},
			{
				name:    "should return an error when code is whitespace",
				setup:   func(t *testing.T) *payment.Payment { return createValidPayment(t) },
				code:    "   ",
				wantErr: payment.ErrInvalidTransactionCode,
			},
			{
				name:    "should return an error when transaction code is already defined",
				setup:   func(t *testing.T) *payment.Payment { return createPaymentWithCode(t) },
				code:    "TXN-456",
				wantErr: payment.ErrTransactionCodeAlreadyDefined,
			},
			{
				name: "should return an error when payment has already been confirmed",
				setup: func(t *testing.T) *payment.Payment {
					p := createPaymentWithCode(t)
					require.NoError(t, p.ConfirmPayment())
					return p
				},
				code:    "TXN-456",
				wantErr: payment.ErrCannotDefineTransactionCodeAfterCompletion,
			},
			{
				name: "should return an error when payment has already been refused",
				setup: func(t *testing.T) *payment.Payment {
					p := createPaymentWithCode(t)
					require.NoError(t, p.RefusePayment())
					return p
				},
				code:    "TXN-456",
				wantErr: payment.ErrCannotDefineTransactionCodeAfterCompletion,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				p := tt.setup(t)

				err := p.DefineTransactionCode(tt.code)

				assert.ErrorIs(t, err, tt.wantErr)
			})
		}
	})
}

func TestPayment_ConfirmPayment(t *testing.T) {
	t.Run("should successfully confirm payment when transaction code has been defined", func(t *testing.T) {
		p := createValidPayment(t)
		require.NoError(t, p.DefineTransactionCode("TXN-123"))

		err := p.ConfirmPayment()

		require.NoError(t, err)
		assert.Equal(t, payment.StatusAuthorized, p.Status, "status should be StatusAuthorized on success")
		assert.NotNil(t, p.PaidAt, "PaidAt should be set on success")
		assert.NotNil(t, p.UpdatedAt, "UpdatedAt should be set on success")
	})

	t.Run("should return an error when state transition is invalid", func(t *testing.T) {
		tests := []struct {
			name    string
			setup   func(t *testing.T) *payment.Payment
			wantErr error
		}{
			{
				name: "should return an error when payment is not pending",
				setup: func(t *testing.T) *payment.Payment {
					p := createPaymentWithCode(t)
					require.NoError(t, p.ConfirmPayment())
					return p
				},
				wantErr: payment.ErrPaymentNotPending,
			},
			{
				name:    "should return an error when transaction code has not been defined",
				setup:   func(t *testing.T) *payment.Payment { return createValidPayment(t) },
				wantErr: payment.ErrTransactionCodeNotDefined,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				p := tt.setup(t)

				err := p.ConfirmPayment()

				assert.ErrorIs(t, err, tt.wantErr)
			})
		}
	})
}

func TestPayment_RefusePayment(t *testing.T) {
	t.Run("should successfully refuse payment when transaction code has been defined", func(t *testing.T) {
		p := createValidPayment(t)
		require.NoError(t, p.DefineTransactionCode("TXN-123"))

		err := p.RefusePayment()

		require.NoError(t, err)
		assert.Equal(t, payment.StatusRefused, p.Status, "status should be StatusRefused on success")
		assert.Nil(t, p.PaidAt, "PaidAt should remain nil on refusal")
		assert.NotNil(t, p.UpdatedAt, "UpdatedAt should be set on success")
	})

	t.Run("should return an error when state transition is invalid", func(t *testing.T) {
		tests := []struct {
			name    string
			setup   func(t *testing.T) *payment.Payment
			wantErr error
		}{
			{
				name: "should return an error when payment is not pending - already refused",
				setup: func(t *testing.T) *payment.Payment {
					p := createPaymentWithCode(t)
					require.NoError(t, p.RefusePayment())
					return p
				},
				wantErr: payment.ErrPaymentNotPending,
			},
			{
				name: "should return an error when payment is not pending - already confirmed",
				setup: func(t *testing.T) *payment.Payment {
					p := createPaymentWithCode(t)
					require.NoError(t, p.ConfirmPayment())
					return p
				},
				wantErr: payment.ErrPaymentNotPending,
			},
			{
				name:    "should return an error when transaction code has not been defined",
				setup:   func(t *testing.T) *payment.Payment { return createValidPayment(t) },
				wantErr: payment.ErrTransactionCodeNotDefined,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				p := tt.setup(t)

				err := p.RefusePayment()

				assert.ErrorIs(t, err, tt.wantErr)
			})
		}
	})
}
