package order

import (
	"github.com/marcosvieirajr/sales-ddd-hexagonal/kernel"
)

type CancelledEvent struct {
	kernel.Event
	OrderID            string             `json:"order_id"`
	CustomerID         string             `json:"customer_id"`
	PaymentID          *string            `json:"payment_id"`
	Status             Status             `json:"status"`
	CancellationReason CancellationReason `json:"cancellation_reason"`
}
