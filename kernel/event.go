package kernel

import "time"

// Event is the base struct embedded in all domain events.
// It carries a unique ID and the UTC timestamp of when the event occurred.
type Event struct {
	ID           string    `json:"id"`
	DateOccurred time.Time `json:"occurred_at"`
}

// EventID returns the event's unique identifier, satisfying the [DomainEvent] interface.
func (e Event) EventID() string {
	return e.ID
}

// OccurredAt returns the UTC timestamp at which the event was created.
func (e Event) OccurredAt() time.Time {
	return e.DateOccurred
}
