package order

// Status represents the fulfillment lifecycle state of an [Order].
type Status int

const (
	StatusCreated    Status = iota // StatusCreated is the initial state of an order after placement.
	StatusPaid                     // StatusPaid indicates the order payment has been confirmed.
	StatusSeparating               // StatusSeparating indicates the order is being picked and packed.
	StatusShipped                  // StatusShipped indicates the order has been dispatched to the carrier.
	StatusDelivered                // StatusDelivered indicates the order has reached the customer.
	StatusCancelled                // StatusCancelled indicates the order has been cancelled.
)
