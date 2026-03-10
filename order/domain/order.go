package order

import (
	"errors"
	"time"

	"github.com/marcosvieirajr/sales-ddd-hexagonal/kernel"
	"github.com/marcosvieirajr/sales-ddd-hexagonal/kernel/errs"
	"github.com/marcosvieirajr/sales-ddd-hexagonal/kernel/guard"
	"github.com/marcosvieirajr/sales-ddd-hexagonal/order/domain/orderitem"
	"github.com/marcosvieirajr/sales-ddd-hexagonal/order/domain/payment"
)

var (
	ErrInvalidCustomerID      = errs.New("ORDER.INVALID_CUSTOMER_ID", "customer ID cannot be null or whitespace")
	ErrInvalidDeliveryAddress = errs.New("ORDER.INVALID_DELIVERY_ADDRESS", "delivery address cannot be zero")
	ErrOrderNotPending        = errs.New("ORDER.NOT_PENDING", "order must be in pending status to perform this operation")
	ErrItemNotFound           = errs.New("ORDER.ITEM_NOT_FOUND", "item not found in order")
	ErrCannotRemoveLastItem   = errs.New("ORDER.CANNOT_REMOVE_LAST_ITEM", "cannot remove the last item from an order")
	ErrNoItems                = errs.New("ORDER.NO_ITEMS", "order must have at least one item to start payment")
	ErrPaymentAlreadyPending  = errs.New("ORDER.PAYMENT_ALREADY_PENDING", "order already has a pending payment")
	ErrOrderNotPaid           = errs.New("ORDER.NOT_PAID", "order must be in paid status to start separating")
	ErrOrderNotSeparating     = errs.New("ORDER.NOT_SEPARATING", "order must be in separating status to be shipped")
	ErrOrderNotShipped        = errs.New("ORDER.NOT_SHIPPED", "order must be in shipped status to be delivered")
	ErrOrderCannotCancel      = errs.New("ORDER.CANNOT_CANCEL", "order cannot be cancelled in its current status")
)

// Order is the aggregate root of the order bounded context.
// It owns the lifecycle of its associated payment and order items.
type Order struct {
	kernel.AggregateRoot
	ID              string
	CustomerID      string
	DeliveryAddress DeliveryAddress
	TotalAmount     float64
	Status          Status
	Number          string
	UpdatedAt       *time.Time

	// ===== Itens ===== //
	items map[string]*orderitem.OrderItem

	// ===== Payment ====== //
	payments    map[string]*payment.Payment
	lastPayment *payment.Payment
}

// NewOrder is a factory that creates a new pending Order, validating customerID (non-blank)
// and address (non-zero).
func NewOrder(customerID string, address *DeliveryAddress) (*Order, error) {
	if err := errors.Join(
		guard.CheckNotNullOrWhiteSpace(customerID, ErrInvalidCustomerID),
		guard.CheckNotZeroValue(address, ErrInvalidDeliveryAddress),
	); err != nil {
		return nil, err
	}

	return &Order{
		ID:              kernel.NewID().String(),
		CustomerID:      customerID,
		DeliveryAddress: *address,
		TotalAmount:     0,
		Status:          StatusPending,
		Number:          generateNumber(),
		items:           make(map[string]*orderitem.OrderItem),
		payments:        make(map[string]*payment.Payment),
	}, nil
}

// AddItem adds or increases the quantity of a product line item; the order must be pending.
func (o *Order) AddItem(productID, productName string, unitPrice float64, quantity int) error {
	if !o.Status.Equals(StatusPending) {
		return ErrOrderNotPending
	}

	if item, exists := o.items[productID]; exists {
		err := item.AddUnits(quantity)
		if err != nil {
			return err
		}

		o.calculateTotalAmount()
		o.updateTimestamp()
		return nil
	}

	item, err := orderitem.NewOrderItem(productID, productName, unitPrice, quantity)
	if err != nil {
		return err
	}

	o.items[productID] = item
	o.calculateTotalAmount()
	o.updateTimestamp()

	return nil
}

// RemoveItem removes a line item from the order; the order must be pending and at least
// one other item must remain.
func (o *Order) RemoveItem(item *orderitem.OrderItem) error {
	if !o.Status.Equals(StatusPending) {
		return ErrOrderNotPending
	}

	if _, exists := o.items[item.ProductID]; !exists {
		return ErrItemNotFound
	}

	if len(o.items) == 1 {
		return ErrCannotRemoveLastItem
	}

	delete(o.items, item.ProductID)

	o.calculateTotalAmount()
	o.updateTimestamp()
	return nil
}

// UpdateDeliveryAddress replaces the delivery address; the order must be pending and
// the new address must be non-zero.
func (o *Order) UpdateDeliveryAddress(newAddress DeliveryAddress) error {
	if !o.Status.Equals(StatusPending) {
		return ErrOrderNotPending
	}

	if newAddress.IsZero() {
		return ErrInvalidDeliveryAddress
	}

	o.DeliveryAddress = newAddress
	o.updateTimestamp()
	return nil
}

// StartPayment creates a new pending Payment for the order; the order must be pending,
// have items, and have no existing pending payment.
func (o *Order) StartPayment(method payment.Method) (*payment.Payment, error) {
	if !o.Status.Equals(StatusPending) {
		return nil, ErrOrderNotPending
	}

	if len(o.items) == 0 {
		return nil, ErrNoItems
	}

	for _, p := range o.payments {
		if p.Status.Equals(payment.StatusPending) {
			return nil, ErrPaymentAlreadyPending
		}
	}

	newPayment, err := payment.NewPayment(o.ID, o.TotalAmount, method)
	if err != nil {
		return nil, err
	}

	o.payments[newPayment.ID] = newPayment
	o.lastPayment = newPayment
	o.updateTimestamp()
	return newPayment, nil
}

// HandleApprovedPaymentEvent transitions the order to Paid when the identified payment
// is approved.
func (o *Order) HandleApprovedPaymentEvent(paymentID string) error {
	if !o.Status.Equals(StatusPending) {
		return ErrOrderNotPending
	}

	if _, exists := o.payments[paymentID]; !exists {
		return nil
	}

	o.Status = StatusPaid
	o.updateTimestamp()
	return nil
}

// HandleRejectedPaymentEvent transitions the order to Cancelled and raises a CancelledEvent
// when the identified payment is rejected.
func (o *Order) HandleRejectedPaymentEvent(paymentID string) error {
	if !o.Status.Equals(StatusPending) {
		return ErrOrderNotPending
	}

	if _, exists := o.payments[paymentID]; !exists {
		return nil
	}

	o.Status = StatusCancelled
	o.updateTimestamp()

	event := newCancelledEvent(o.ID, o.CustomerID, o.Status, CancellationReasonPaymentError, paymentID)
	o.AddDomainEvent(event)
	return nil
}

// MarkAsSeparating advances the order to the Separating status; the order must be Paid.
func (o *Order) MarkAsSeparating() error {
	if !o.Status.Equals(StatusPaid) {
		return ErrOrderNotPaid
	}

	o.Status = StatusSeparating
	o.updateTimestamp()
	return nil
}

// MarkAsShipped advances the order to the Shipped status and raises a ShippedEvent;
// the order must be Separating.
func (o *Order) MarkAsShipped() error {
	if !o.Status.Equals(StatusSeparating) {
		return ErrOrderNotSeparating
	}

	o.Status = StatusShipped
	o.updateTimestamp()

	event := newShippedEvent(o.ID, o.CustomerID, o.DeliveryAddress)
	o.AddDomainEvent(event)
	return nil
}

// MarkAsDelivered advances the order to the Delivered status and raises a DeliveredEvent;
// the order must be Shipped.
func (o *Order) MarkAsDelivered() error {
	if !o.Status.Equals(StatusShipped) {
		return ErrOrderNotShipped
	}

	o.Status = StatusDelivered
	o.updateTimestamp()

	event := newDeliveredEvent(o.ID, o.CustomerID)
	o.AddDomainEvent(event)
	return nil
}

// Cancel cancels the order and raises a CancelledEvent; the order must be in a
// cancellable status.
func (o *Order) Cancel(reason CancellationReason) error {
	if !o.Status.Equals(StatusShipped) &&
		!o.Status.Equals(StatusDelivered) {
		return ErrOrderCannotCancel
	}

	o.Status = StatusCancelled
	o.updateTimestamp()

	var paymentID string
	if o.lastPayment != nil {
		paymentID = o.lastPayment.ID
	}

	event := newCancelledEvent(o.ID, o.CustomerID, o.Status, reason, paymentID)
	o.AddDomainEvent(event)
	return nil
}

func (o *Order) updateTimestamp() {
	timestamp := time.Now().UTC()
	o.UpdatedAt = &timestamp
}

func (o *Order) calculateTotalAmount() {
	totalAmount := 0.0
	for _, item := range o.items {
		totalAmount = +totalAmount + item.TotalPrice
	}
	o.TotalAmount = totalAmount
}

func generateNumber() string {
	return "PED-" + kernel.NewID().String()[:8] // TODO: reimplement
}
