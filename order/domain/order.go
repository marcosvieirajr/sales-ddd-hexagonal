package order

// Order is the aggregate root of the order bounded context.
// It owns the lifecycle of its associated payment and order items.
type Order struct {
	ID         string
	Status     Status
	TotalPrice float64
}
