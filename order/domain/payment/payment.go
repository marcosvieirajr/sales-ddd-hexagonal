package payment

import (
	"errors"
	"time"

	"github.com/marcosvieirajr/sales-ddd-hexagonal/kernel"
	"github.com/marcosvieirajr/sales-ddd-hexagonal/kernel/errs"
	"github.com/marcosvieirajr/sales-ddd-hexagonal/kernel/guard"
)

var (
	ErrInvalidOrderID                             = errs.New("PAYMENT.INVALID_ORDER_ID", "order ID cannot be null or whitespace")
	ErrInvalidPaymentAmount                       = errs.New("PAYMENT.INVALID_AMOUNT", "payment amount must be greater than zero")
	ErrInvalidTransactionCode                     = errs.New("PAYMENT.INVALID_TRANSACTION_CODE", "transaction code cannot be null or whitespace")
	ErrTransactionCodeAlreadyDefined              = errs.New("PAYMENT.TRANSACTION_CODE_ALREADY_DEFINED", "transaction code has already been defined")
	ErrCannotDefineTransactionCodeAfterCompletion = errs.New("PAYMENT.TRANSACTION_CODE_AFTER_COMPLETION", "transaction code cannot be defined after payment has been confirmed or refused")
	ErrPaymentNotPending                          = errs.New("PAYMENT.NOT_PENDING", "payment is not in pending status")
	ErrTransactionCodeNotDefined                  = errs.New("PAYMENT.TRANSACTION_CODE_NOT_DEFINED", "transaction code has not been defined yet")
)

// Payment is an entity of the Order aggregate that represents a payment transaction.
// It is created in [StatusPending] and transitions to [StatusAuthorized] or [StatusRefused]
// via [ConfirmPayment] or [RefusePayment] respectively, after a transaction code has been
// assigned with [DefineTransactionCode].
type Payment struct {
	ID              string
	OrderID         string
	Amount          float64 // TODO: create a value object using a more precise type for money
	Method          Method
	Status          Status
	PaidAt          *time.Time
	UpdatedAt       *time.Time
	TransactionCode *string
}

// NewPayment creates a new [Payment] for the given order with the specified amount and payment method.
// orderID must be non-empty and non-whitespace; amount must be strictly positive.
// The payment is initialized in [StatusPending] with no transaction code assigned.
//
// If multiple fields are invalid, all violations are collected and returned as a
// single joined error, allowing callers to inspect every failure via [errors.Is].
func NewPayment(orderID string, amount float64, method Method) (*Payment, error) {
	// the order ID cannot be null or whitespace, and the amount must be greater than zero.
	if err := errors.Join(
		guard.CheckNotNullOrWhiteSpace(orderID, ErrInvalidOrderID),
		guard.CheckNotZeroOrNegative(amount, ErrInvalidPaymentAmount),
	); err != nil {
		return nil, err
	}

	return &Payment{
		ID:      kernel.GenerateID(),
		OrderID: orderID,
		Method:  method,
		Status:  StatusPending,
		Amount:  amount,
	}, nil
}

// ConfirmPayment transitions the payment from [StatusPending] to [StatusAuthorized],
// recording the current UTC time as PaidAt and refreshing UpdatedAt.
// Returns [ErrPaymentNotPending] if the payment is not pending, or
// [ErrTransactionCodeNotDefined] if no transaction code has been set.
func (p *Payment) ConfirmPayment() error {
	// the payment can only be confirmed if it is currently pending and has a transaction code defined.
	if err := errors.Join(
		p.checkStatusEqual(StatusPending, ErrPaymentNotPending),
		guard.CheckNotNil(p.TransactionCode, ErrTransactionCodeNotDefined),
	); err != nil {
		return err
	}

	now := time.Now().UTC()
	p.PaidAt = &now
	p.Status = StatusAuthorized
	p.updateTimestamp()
	p.AddDomainEvent(ApprovedEvent{}) // TODO: add more details to the event (e.g. order ID, amount, etc.) and test that it is emitted correctly.

	return nil
}

// RefusePayment transitions the payment from [StatusPending] to [StatusRefused],
// refreshing UpdatedAt.
// Returns [ErrPaymentNotPending] if the payment is not pending, or
// [ErrTransactionCodeNotDefined] if no transaction code has been set.
func (p *Payment) RefusePayment() error {
	// the payment can only be refused if it is currently pending and has a transaction code defined.
	if err := errors.Join(
		p.checkStatusEqual(StatusPending, ErrPaymentNotPending),
		guard.CheckNotNil(p.TransactionCode, ErrTransactionCodeNotDefined),
	); err != nil {
		return err
	}

	p.Status = StatusRefused
	p.updateTimestamp()
	p.AddDomainEvent(RefusedEvent{}) // TODO: add more details to the event (e.g. order ID, amount, etc.) and test that it is emitted correctly.

	return nil
}

// DefineTransactionCode assigns the external transaction code returned by the payment gateway.
// code must be non-empty and non-whitespace.
// Returns [ErrCannotDefineTransactionCodeAfterCompletion] if the payment is no longer pending,
// [ErrTransactionCodeAlreadyDefined] if a code has already been set, or
// [ErrInvalidTransactionCode] if code is blank.
func (p *Payment) DefineTransactionCode(code string) error {
	// validate that the code is not null or whitespace, that no code has been defined yet,
	// and that the payment is pending (i.e. not already approved or refused).
	if err := errors.Join(
		p.checkStatusEqual(StatusPending, ErrCannotDefineTransactionCodeAfterCompletion),
		guard.CheckNotNullOrWhiteSpace(code, ErrInvalidTransactionCode),
		guard.CheckNil(p.TransactionCode, ErrTransactionCodeAlreadyDefined),
	); err != nil {
		return err
	}

	p.TransactionCode = &code
	p.updateTimestamp()

	return nil
}

type DomainEvent interface {
	OccurredAt() time.Time
}

func (p *Payment) AddDomainEvent(event DomainEvent) {
	// TODO: implement and test...
}

func (p *Payment) updateTimestamp() {
	timestamp := time.Now().UTC()
	p.UpdatedAt = &timestamp
}

func (p *Payment) checkStatusEqual(other Status, err error) error {
	if !p.Status.Equals(other) {
		return err
	}
	return nil
}

func (p *Payment) generateTransactionCode() {
	if p.TransactionCode != nil {
		return
	}

	c := `LOCAL-` + kernel.GenerateID()
	_ = p.DefineTransactionCode(c)
}
