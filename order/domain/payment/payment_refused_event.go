package payment

import "time"

type Event struct {
	DateOccurred time.Time `json:"occurred_at"`
}

func (e Event) OccurredAt() time.Time {
	return e.DateOccurred
}

// RefusedEvent represents the event when a payment is refused.
type RefusedEvent struct {
	Event
	PaymentID       string  `json:"payment_id"`
	OrderID         string  `json:"order_id"`
	Amount          float64 `json:"amount"`
	TransactionCode *string `json:"transaction_code"`
}

func NewRefusedEvent(paymentID, orderID string, amount float64, transactionCode *string) RefusedEvent {
	return RefusedEvent{
		Event: Event{
			DateOccurred: time.Now().UTC(),
		},
		PaymentID:       paymentID,
		OrderID:         orderID,
		Amount:          amount,
		TransactionCode: transactionCode,
	}
}
