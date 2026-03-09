package order

import (
	"github.com/marcosvieirajr/sales-ddd-hexagonal/kernel"
)

type DeliveredEvent struct {
	kernel.Event
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
}
