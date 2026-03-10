package order

import (
	"time"

	"github.com/marcosvieirajr/sales-ddd-hexagonal/kernel"
)

// ShippedEvent is a domain event raised when an Order is dispatched,
// carrying the delivery address.
type ShippedEvent struct {
	kernel.Event
	OrderID         string          `json:"order_id"`
	CustomerID      string          `json:"customer_id"`
	DeliveryAddress DeliveryAddress `json:"delivery_address"`
}

func newShippedEvent(orderID string, customerID string, deliveryAddress DeliveryAddress) *ShippedEvent {
	return &ShippedEvent{
		Event: kernel.Event{
			ID:           kernel.NewID().String(),
			DateOccurred: time.Now().UTC(),
		},
		OrderID:         orderID,
		CustomerID:      customerID,
		DeliveryAddress: deliveryAddress,
	}
}
