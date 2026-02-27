package payment

import "time"

// ApprovedEvent represents the event when a payment is approved.
type ApprovedEvent struct {
	Event
	PaymentID       string  `json:"payment_id"`
	OrderID         string  `json:"order_id"`
	Amount          float64 `json:"amount"`
	TransactionCode *string `json:"transaction_code"`
}

func NewApprovedEvent(paymentID, orderID string, amount float64, transactionCode *string) RefusedEvent {
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
