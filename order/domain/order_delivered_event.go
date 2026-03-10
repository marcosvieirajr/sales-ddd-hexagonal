package order

import (
	"time"

	"github.com/marcosvieirajr/sales-ddd-hexagonal/kernel"
)

// DeliveredEvent is a domain event raised when an Order is successfully delivered
// to the customer.
type DeliveredEvent struct {
	kernel.Event
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
}

func newDeliveredEvent(orderID string, customerID string) *DeliveredEvent {
	return &DeliveredEvent{
		Event: kernel.Event{
			ID:           kernel.NewID().String(),
			DateOccurred: time.Now().UTC(),
		},
		OrderID:    orderID,
		CustomerID: customerID,
	}
}
