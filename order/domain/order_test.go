package order_test

import (
	"testing"

	"github.com/marcosvieirajr/sales-ddd-hexagonal/kernel"
	order "github.com/marcosvieirajr/sales-ddd-hexagonal/order/domain"
	"github.com/marcosvieirajr/sales-ddd-hexagonal/order/domain/orderitem"
	"github.com/marcosvieirajr/sales-ddd-hexagonal/order/domain/payment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== Helpers ==================== //

func createValidAddress(t *testing.T) *order.DeliveryAddress {
	t.Helper()
	return kernel.Must(order.NewDeliveryAddress("12345-678", "Rua das Flores", "100", "", "Centro", "São Paulo", "SP", "Brasil"))
}

func createValidOrder(t *testing.T) *order.Order {
	t.Helper()
	return kernel.Must(order.NewOrder("cust-123", createValidAddress(t)))
}

func createOrderWithItems(t *testing.T) *order.Order {
	t.Helper()
	o := createValidOrder(t)
	require.NoError(t, o.AddItem("prod-1", "Widget", 50.0, 2))
	return o
}

func driveOrderToPaid(t *testing.T) *order.Order {
	t.Helper()
	o := createOrderWithItems(t)
	p, err := o.StartPayment(payment.MethodCreditCard)
	require.NoError(t, err)
	require.NoError(t, o.HandleApprovedPaymentEvent(p.ID))
	return o
}

func driveOrderToSeparating(t *testing.T) *order.Order {
	t.Helper()
	o := driveOrderToPaid(t)
	require.NoError(t, o.MarkAsSeparating())
	return o
}

func driveOrderToShipped(t *testing.T) *order.Order {
	t.Helper()
	o := driveOrderToSeparating(t)
	require.NoError(t, o.MarkAsShipped())
	return o
}

func driveOrderToDelivered(t *testing.T) *order.Order {
	t.Helper()
	o := driveOrderToShipped(t)
	require.NoError(t, o.MarkAsDelivered())
	return o
}

// ==================== Tests ==================== //

func TestNewOrder(t *testing.T) {
	t.Run("should successfully create a new order with valid input", func(t *testing.T) {
		addr := createValidAddress(t)

		got, err := order.NewOrder("cust-123", addr)

		require.NoError(t, err)
		require.NotNil(t, got)
		assert.NotEmpty(t, got.ID, "ID should be generated")
		assert.Equal(t, "cust-123", got.CustomerID)
		assert.Equal(t, order.StatusPending, got.Status, "status should be Pending")
		assert.Equal(t, 0.0, got.TotalAmount, "TotalAmount should be zero on creation")
		assert.Nil(t, got.UpdatedAt, "UpdatedAt should be nil on creation")
	})

	t.Run("should return an error when input is invalid", func(t *testing.T) {
		addr := createValidAddress(t)
		tests := []struct {
			name       string
			customerID string
			address    *order.DeliveryAddress
			wantErr    error
		}{
			{
				name:       "should return an error when customerID is empty",
				customerID: "",
				address:    addr,
				wantErr:    order.ErrInvalidCustomerID,
			},
			{
				name:       "should return an error when customerID is whitespace",
				customerID: "   ",
				address:    addr,
				wantErr:    order.ErrInvalidCustomerID,
			},
			{
				name:       "should return an error when address is nil",
				customerID: "cust-123",
				address:    nil,
				wantErr:    order.ErrInvalidDeliveryAddress,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := order.NewOrder(tt.customerID, tt.address)

				assert.Nil(t, got)
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			})
		}
	})
}

func TestOrder_AddItem(t *testing.T) {
	t.Run("should successfully add a new item and update TotalAmount", func(t *testing.T) {
		o := createValidOrder(t)

		err := o.AddItem("prod-1", "Widget", 50.0, 2)

		require.NoError(t, err)
		assert.Equal(t, 100.0, o.TotalAmount, "TotalAmount should be 50 * 2 = 100")
		assert.NotNil(t, o.UpdatedAt, "UpdatedAt should be set on success")
	})

	t.Run("should successfully increase quantity when item already exists", func(t *testing.T) {
		o := createValidOrder(t)
		require.NoError(t, o.AddItem("prod-1", "Widget", 50.0, 2))

		err := o.AddItem("prod-1", "Widget", 50.0, 3)

		require.NoError(t, err)
		assert.Equal(t, 250.0, o.TotalAmount, "TotalAmount should be 50 * 5 = 250")
	})

	t.Run("should return an error when order is not pending", func(t *testing.T) {
		tests := []struct {
			name  string
			setup func(t *testing.T) *order.Order
		}{
			{name: "status Paid", setup: driveOrderToPaid},
			{name: "status Separating", setup: driveOrderToSeparating},
			{name: "status Shipped", setup: driveOrderToShipped},
			{name: "status Delivered", setup: driveOrderToDelivered},
			{
				name: "status Cancelled",
				setup: func(t *testing.T) *order.Order {
					o := driveOrderToShipped(t)
					require.NoError(t, o.Cancel(order.CancellationReasonCustomerCancelled))
					return o
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				o := tt.setup(t)

				err := o.AddItem("prod-2", "Gadget", 10.0, 1)

				assert.ErrorIs(t, err, order.ErrOrderNotPending)
			})
		}
	})
}

func TestOrder_RemoveItem(t *testing.T) {
	t.Run("should successfully remove an existing item and recalculate TotalAmount", func(t *testing.T) {
		o := createValidOrder(t)
		require.NoError(t, o.AddItem("prod-1", "Widget", 50.0, 2))
		require.NoError(t, o.AddItem("prod-2", "Gadget", 10.0, 1))
		itemRef := kernel.Must(orderitem.NewOrderItem("prod-2", "Gadget", 10.0, 1))

		err := o.RemoveItem(itemRef)

		require.NoError(t, err)
		assert.Equal(t, 100.0, o.TotalAmount, "TotalAmount should be 50*2=100 after removing prod-2")
		assert.NotNil(t, o.UpdatedAt, "UpdatedAt should be set on success")
	})

	t.Run("should return an error when order is not pending", func(t *testing.T) {
		o := driveOrderToPaid(t)
		item := kernel.Must(orderitem.NewOrderItem("prod-1", "Widget", 50.0, 2))

		err := o.RemoveItem(item)

		assert.ErrorIs(t, err, order.ErrOrderNotPending)
	})

	t.Run("should return an error when item is not in the order", func(t *testing.T) {
		o := createOrderWithItems(t)
		unknownItem := kernel.Must(orderitem.NewOrderItem("prod-unknown", "Unknown", 5.0, 1))

		err := o.RemoveItem(unknownItem)

		assert.ErrorIs(t, err, order.ErrItemNotFound)
	})

	t.Run("should return an error when attempting to remove the last item", func(t *testing.T) {
		o := createValidOrder(t)
		require.NoError(t, o.AddItem("prod-1", "Widget", 50.0, 2))
		item := kernel.Must(orderitem.NewOrderItem("prod-1", "Widget", 50.0, 2))

		err := o.RemoveItem(item)

		assert.ErrorIs(t, err, order.ErrCannotRemoveLastItem)
	})
}

func TestOrder_UpdateDeliveryAddress(t *testing.T) {
	t.Run("should successfully update delivery address", func(t *testing.T) {
		o := createValidOrder(t)
		newAddr := kernel.Must(order.NewDeliveryAddress("98765-432", "Av. Brasil", "500", "Apto 1", "Jardins", "Rio de Janeiro", "RJ", "Brasil"))

		err := o.UpdateDeliveryAddress(*newAddr)

		require.NoError(t, err)
		assert.True(t, o.DeliveryAddress.Equals(newAddr), "DeliveryAddress should be replaced")
		assert.NotNil(t, o.UpdatedAt, "UpdatedAt should be set on success")
	})

	t.Run("should return an error when order is not pending", func(t *testing.T) {
		o := driveOrderToPaid(t)
		newAddr := kernel.Must(order.NewDeliveryAddress("98765-432", "Av. Brasil", "500", "", "Jardins", "Rio de Janeiro", "RJ", "Brasil"))

		err := o.UpdateDeliveryAddress(*newAddr)

		assert.ErrorIs(t, err, order.ErrOrderNotPending)
	})

	t.Run("should return an error when address is zero value", func(t *testing.T) {
		o := createValidOrder(t)

		err := o.UpdateDeliveryAddress(order.DeliveryAddress{})

		assert.ErrorIs(t, err, order.ErrInvalidDeliveryAddress)
	})
}

func TestOrder_StartPayment(t *testing.T) {
	t.Run("should successfully start a payment and store it", func(t *testing.T) {
		o := createOrderWithItems(t)

		p, err := o.StartPayment(payment.MethodCreditCard)

		require.NoError(t, err)
		require.NotNil(t, p)
		assert.NotEmpty(t, p.ID)
		assert.Equal(t, payment.StatusPending, p.Status, "payment status should be Pending")
		assert.NotNil(t, o.UpdatedAt, "UpdatedAt should be set on success")
	})

	t.Run("should return an error when order is not pending", func(t *testing.T) {
		o := driveOrderToPaid(t)

		p, err := o.StartPayment(payment.MethodCreditCard)

		assert.Nil(t, p)
		assert.ErrorIs(t, err, order.ErrOrderNotPending)
	})

	t.Run("should return an error when order has no items", func(t *testing.T) {
		o := createValidOrder(t)

		p, err := o.StartPayment(payment.MethodCreditCard)

		assert.Nil(t, p)
		assert.ErrorIs(t, err, order.ErrNoItems)
	})

	t.Run("should return an error when a pending payment already exists", func(t *testing.T) {
		o := createOrderWithItems(t)
		_, err := o.StartPayment(payment.MethodCreditCard)
		require.NoError(t, err)

		p2, err := o.StartPayment(payment.MethodCreditCard)

		assert.Nil(t, p2)
		assert.ErrorIs(t, err, order.ErrPaymentAlreadyPending)
	})
}

func TestOrder_HandleApprovedPaymentEvent(t *testing.T) {
	t.Run("should transition order to Paid when payment is approved", func(t *testing.T) {
		o := createOrderWithItems(t)
		p, err := o.StartPayment(payment.MethodCreditCard)
		require.NoError(t, err)

		err = o.HandleApprovedPaymentEvent(p.ID)

		require.NoError(t, err)
		assert.Equal(t, order.StatusPaid, o.Status, "status should be Paid")
		assert.NotNil(t, o.UpdatedAt, "UpdatedAt should be set on success")
	})

	t.Run("should return an error when order is not pending", func(t *testing.T) {
		o := driveOrderToPaid(t)

		err := o.HandleApprovedPaymentEvent("any-id")

		assert.ErrorIs(t, err, order.ErrOrderNotPending)
	})

	t.Run("should be a no-op when paymentID is unknown", func(t *testing.T) {
		o := createOrderWithItems(t)

		err := o.HandleApprovedPaymentEvent("unknown-payment-id")

		require.NoError(t, err)
		assert.Equal(t, order.StatusPending, o.Status, "status should remain Pending")
	})
}

func TestOrder_HandleRejectedPaymentEvent(t *testing.T) {
	t.Run("should transition order to Cancelled when payment is rejected", func(t *testing.T) {
		o := createOrderWithItems(t)
		p, err := o.StartPayment(payment.MethodCreditCard)
		require.NoError(t, err)

		err = o.HandleRejectedPaymentEvent(p.ID)

		require.NoError(t, err)
		assert.Equal(t, order.StatusCancelled, o.Status, "status should be Cancelled")
		assert.NotNil(t, o.UpdatedAt, "UpdatedAt should be set on success")
	})

	t.Run("should return an error when order is not pending", func(t *testing.T) {
		o := driveOrderToPaid(t)

		err := o.HandleRejectedPaymentEvent("any-id")

		assert.ErrorIs(t, err, order.ErrOrderNotPending)
	})

	t.Run("should be a no-op when paymentID is unknown", func(t *testing.T) {
		o := createOrderWithItems(t)

		err := o.HandleRejectedPaymentEvent("unknown-payment-id")

		require.NoError(t, err)
		assert.Equal(t, order.StatusPending, o.Status, "status should remain Pending")
	})
}

func TestOrder_MarkAsSeparating(t *testing.T) {
	t.Run("should transition order from Paid to Separating", func(t *testing.T) {
		o := driveOrderToPaid(t)

		err := o.MarkAsSeparating()

		require.NoError(t, err)
		assert.Equal(t, order.StatusSeparating, o.Status, "status should be Separating")
		assert.NotNil(t, o.UpdatedAt, "UpdatedAt should be set on success")
	})

	t.Run("should return an error when order is not Paid", func(t *testing.T) {
		tests := []struct {
			name  string
			setup func(t *testing.T) *order.Order
		}{
			{name: "status Pending", setup: createValidOrder},
			{name: "status Separating", setup: driveOrderToSeparating},
			{name: "status Shipped", setup: driveOrderToShipped},
			{name: "status Delivered", setup: driveOrderToDelivered},
			{
				name: "status Cancelled",
				setup: func(t *testing.T) *order.Order {
					o := driveOrderToShipped(t)
					require.NoError(t, o.Cancel(order.CancellationReasonCustomerCancelled))
					return o
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				o := tt.setup(t)

				err := o.MarkAsSeparating()

				assert.ErrorIs(t, err, order.ErrOrderNotPaid)
			})
		}
	})
}

func TestOrder_MarkAsShipped(t *testing.T) {
	t.Run("should transition order from Separating to Shipped", func(t *testing.T) {
		o := driveOrderToSeparating(t)

		err := o.MarkAsShipped()

		require.NoError(t, err)
		assert.Equal(t, order.StatusShipped, o.Status, "status should be Shipped")
		assert.NotNil(t, o.UpdatedAt, "UpdatedAt should be set on success")
	})

	t.Run("should return an error when order is not Separating", func(t *testing.T) {
		tests := []struct {
			name  string
			setup func(t *testing.T) *order.Order
		}{
			{name: "status Pending", setup: createValidOrder},
			{name: "status Paid", setup: driveOrderToPaid},
			{name: "status Shipped", setup: driveOrderToShipped},
			{name: "status Delivered", setup: driveOrderToDelivered},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				o := tt.setup(t)

				err := o.MarkAsShipped()

				assert.ErrorIs(t, err, order.ErrOrderNotSeparating)
			})
		}
	})
}

func TestOrder_MarkAsDelivered(t *testing.T) {
	t.Run("should transition order from Shipped to Delivered", func(t *testing.T) {
		o := driveOrderToShipped(t)

		err := o.MarkAsDelivered()

		require.NoError(t, err)
		assert.Equal(t, order.StatusDelivered, o.Status, "status should be Delivered")
		assert.NotNil(t, o.UpdatedAt, "UpdatedAt should be set on success")
	})

	t.Run("should return an error when order is not Shipped", func(t *testing.T) {
		tests := []struct {
			name  string
			setup func(t *testing.T) *order.Order
		}{
			{name: "status Pending", setup: createValidOrder},
			{name: "status Paid", setup: driveOrderToPaid},
			{name: "status Separating", setup: driveOrderToSeparating},
			{name: "status Delivered", setup: driveOrderToDelivered},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				o := tt.setup(t)

				err := o.MarkAsDelivered()

				assert.ErrorIs(t, err, order.ErrOrderNotShipped)
			})
		}
	})
}

func TestOrder_Cancel(t *testing.T) {
	t.Run("should successfully cancel from Shipped", func(t *testing.T) {
		o := driveOrderToShipped(t)

		err := o.Cancel(order.CancellationReasonCustomerCancelled)

		require.NoError(t, err)
		assert.Equal(t, order.StatusCancelled, o.Status, "status should be Cancelled")
		assert.NotNil(t, o.UpdatedAt, "UpdatedAt should be set on success")
	})

	t.Run("should successfully cancel from Delivered", func(t *testing.T) {
		o := driveOrderToDelivered(t)

		err := o.Cancel(order.CancellationReasonCustomerCancelled)

		require.NoError(t, err)
		assert.Equal(t, order.StatusCancelled, o.Status, "status should be Cancelled")
		assert.NotNil(t, o.UpdatedAt, "UpdatedAt should be set on success")
	})

	t.Run("should return an error when order cannot be cancelled", func(t *testing.T) {
		tests := []struct {
			name  string
			setup func(t *testing.T) *order.Order
		}{
			{name: "status Pending", setup: createValidOrder},
			{name: "status Paid", setup: driveOrderToPaid},
			{name: "status Separating", setup: driveOrderToSeparating},
			{
				name: "status Cancelled",
				setup: func(t *testing.T) *order.Order {
					o := driveOrderToShipped(t)
					require.NoError(t, o.Cancel(order.CancellationReasonCustomerCancelled))
					return o
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				o := tt.setup(t)

				err := o.Cancel(order.CancellationReasonCustomerCancelled)

				assert.ErrorIs(t, err, order.ErrOrderCannotCancel)
			})
		}
	})
}
