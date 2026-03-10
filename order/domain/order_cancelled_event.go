package order

import (
	"strings"
	"time"

	"github.com/marcosvieirajr/sales-ddd-hexagonal/kernel"
)

// CancelledEvent is a domain event raised when an Order is cancelled,
// carrying the cancellation reason and optional payment ID.
type CancelledEvent struct {
	kernel.Event
	OrderID            string             `json:"order_id"`
	CustomerID         string             `json:"customer_id"`
	PaymentID          *string            `json:"payment_id"`
	Status             Status             `json:"status"`
	CancellationReason CancellationReason `json:"cancellation_reason"`
}

func newCancelledEvent(orderID string, customerID string, status Status, reason CancellationReason, paymentID string) *CancelledEvent {
	e := CancelledEvent{
		Event: kernel.Event{
			ID:           kernel.NewID().String(),
			DateOccurred: time.Now().UTC(),
		},
		OrderID:            orderID,
		CustomerID:         customerID,
		Status:             status,
		CancellationReason: reason,
	}

	if strings.TrimSpace(paymentID) != "" {
		e.PaymentID = &paymentID
	}

	return &e
}
