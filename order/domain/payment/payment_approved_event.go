package payment

import (
	"time"

	"github.com/marcosvieirajr/sales-ddd-hexagonal/kernel"
)

// ApprovedEvent represents the event when a payment is approved.
type ApprovedEvent struct {
	kernel.Event
	PaymentID       string  `json:"payment_id"`
	OrderID         string  `json:"order_id"`
	Amount          float64 `json:"amount"`
	TransactionCode *string `json:"transaction_code"`
}

// NewApprovedEvent constructs an ApprovedEvent with the current UTC timestamp.
func NewApprovedEvent(paymentID, orderID string, amount float64, transactionCode *string) RefusedEvent {
	return RefusedEvent{
		Event: kernel.Event{
			ID:           kernel.NewID().String(),
			DateOccurred: time.Now().UTC(),
		},
		PaymentID:       paymentID,
		OrderID:         orderID,
		Amount:          amount,
		TransactionCode: transactionCode,
	}
}
