package orderitem

import (
	"errors"
	"time"

	"github.com/marcosvieirajr/sales-ddd-hexagonal/shared"
	"github.com/marcosvieirajr/sales-ddd-hexagonal/shared/errs"
)

var (
	ErrInvalidProductID         = errs.New("ORDER_ITEM.INVALID_PRODUCT_ID", "product ID cannot be null or whitespace")
	ErrInvalidProductName       = errs.New("ORDER_ITEM.INVALID_PRODUCT_NAME", "product name cannot be null or whitespace")
	ErrInvalidUnitPrice         = errs.New("ORDER_ITEM.INVALID_UNIT_PRICE", "unit price must be greater than zero")
	ErrInvalidQuantity          = errs.New("ORDER_ITEM.INVALID_QUANTITY", "quantity must be greater than zero")
	ErrNegativeDiscount         = errs.New("ORDER_ITEM.NEGATIVE_DISCOUNT", "discount cannot be negative")
	ErrDiscountExceedsUnitPrice = errs.New("ORDER_ITEM.DISCOUNT_EXCEEDS_PRICE", "discount cannot be greater than unit price")
	ErrInvalidUnits             = errs.New("ORDER_ITEM.INVALID_UNITS", "units cannot be zero or negative")
	ErrInsufficientQuantity     = errs.New("ORDER_ITEM.INSUFFICIENT_QUANTITY", "units to remove cannot be greater than or equal to current quantity")
)

// OrderItem is an entity of the Order aggregate that represents a single line item
// within an order, associating a product with a quantity, unit price, and optional
// discount. TotalPrice is automatically maintained as (UnitPrice × Quantity) − DiscountApplied.
type OrderItem struct {
	ID              string
	ProductID       string
	ProductName     string
	UnitPrice       float64
	Quantity        int
	DiscountApplied float64
	TotalPrice      float64
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}

// NewOrderItem constructs and validates a new [OrderItem] for the given product.
// productID and productName must be non-empty and non-whitespace; unitPrice and
// quantity must be strictly positive. DiscountApplied is initialized to zero and
// TotalPrice is computed immediately.
//
// If multiple fields are invalid, all violations are collected and returned as a
// single joined error, allowing callers to inspect every failure via [errors.Is].
func NewOrderItem(productID, productName string, unitPrice float64, quantity int) (*OrderItem, error) {
	if err := errors.Join(
		shared.CheckNotNullOrWhiteSpace(productID, ErrInvalidProductID),
		shared.CheckNotNullOrWhiteSpace(productName, ErrInvalidProductName),
		shared.CheckNotZeroOrNegative(unitPrice, ErrInvalidUnitPrice),
		shared.CheckNotZeroOrNegative(float64(quantity), ErrInvalidQuantity),
	); err != nil {
		return nil, err
	}

	oi := OrderItem{
		ID:          shared.GenerateID(),
		ProductID:   productID,
		ProductName: productName,
		UnitPrice:   unitPrice,
		Quantity:    quantity,
		CreatedAt:   time.Now().UTC(),
	}

	oi.calculateTotalPrice()

	return &oi, nil
}

// ApplyDiscount sets the discount applied to this item's unit price.
// discount must be non-negative and must not exceed [OrderItem.UnitPrice].
// TotalPrice is recalculated after a successful update.
func (oi *OrderItem) ApplyDiscount(discount float64) error {
	if discount < 0 {
		return ErrNegativeDiscount
	}
	if discount > oi.UnitPrice {
		return ErrDiscountExceedsUnitPrice
	}

	oi.DiscountApplied = discount
	oi.calculateTotalPrice()
	oi.updateTimestamp()

	return nil
}

// AddUnits increases the item quantity by units, which must be strictly positive.
// units must be strictly positive.
// TotalPrice is recalculated after a successful update.
func (oi *OrderItem) AddUnits(units int) error {
	// the units to add must be greater than zero.
	if units <= 0 {
		return ErrInvalidUnits
	}

	oi.Quantity += units
	oi.calculateTotalPrice()
	oi.updateTimestamp()

	return nil
}

// RemoveUnits decreases the item quantity by units.
// units must be strictly positive and less than the current quantity
// (at least one unit must remain). TotalPrice is recalculated after a successful update.
func (oi *OrderItem) RemoveUnits(units int) error {
	// the units to remove must be greater than zero and less than the current quantity.
	if units <= 0 {
		return ErrInvalidUnits
	}
	if units >= oi.Quantity {
		return ErrInsufficientQuantity
	}

	oi.Quantity -= units
	oi.calculateTotalPrice()
	oi.updateTimestamp()

	return nil
}

// UpdateUnitPrice sets a new unit price for the item.
// value must be strictly positive. TotalPrice is recalculated after a successful update.
func (oi *OrderItem) UpdateUnitPrice(value float64) error {
	// the unit price must be greater than zero.
	if value <= 0 {
		return ErrInvalidUnitPrice
	}

	oi.UnitPrice = value
	oi.calculateTotalPrice()
	oi.updateTimestamp()

	return nil
}

// Equals reports whether oi and other represent the same order item by comparing IDs.
// It returns false if other is nil.
func (oi *OrderItem) Equals(other *OrderItem) bool {
	if other == nil {
		return false
	}
	return oi.ID == other.ID
}

func (oi *OrderItem) calculateTotalPrice() {
	oi.TotalPrice = (oi.UnitPrice * float64(oi.Quantity)) - oi.DiscountApplied
}

func (oi *OrderItem) updateTimestamp() {
	timestamp := time.Now().UTC()
	oi.UpdatedAt = &timestamp
}
