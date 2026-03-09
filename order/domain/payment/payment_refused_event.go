package payment

import (
	"time"

	"github.com/marcosvieirajr/sales-ddd-hexagonal/kernel"
)

// RefusedEvent represents the event when a payment is refused.
type RefusedEvent struct {
	kernel.Event
	PaymentID       string  `json:"payment_id"`
	OrderID         string  `json:"order_id"`
	Amount          float64 `json:"amount"`
	TransactionCode *string `json:"transaction_code"`
}

func NewRefusedEvent(paymentID, orderID string, amount float64, transactionCode *string) RefusedEvent {
	return RefusedEvent{
		Event: kernel.Event{
			DateOccurred: time.Now().UTC(),
		},
		PaymentID:       paymentID,
		OrderID:         orderID,
		Amount:          amount,
		TransactionCode: transactionCode,
	}
}
