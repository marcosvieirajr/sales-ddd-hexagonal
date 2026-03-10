package kernel

import "time"

// DomainEvent is the interface that all domain events must implement.
// EventID returns a unique event identifier used for deduplication in [AggregateRoot].
type DomainEvent interface {
	EventID() string
	OccurredAt() time.Time
}

// AggregateRoot is an embeddable struct that manages the collection of domain events
// raised by an aggregate. Embed it in any aggregate root to gain event-sourcing support.
type AggregateRoot struct {
	events map[string]DomainEvent
}

// AddDomainEvent registers a domain event, keyed by its EventID to prevent duplicates.
func (o *AggregateRoot) AddDomainEvent(event DomainEvent) {
	if o.events == nil {
		o.events = make(map[string]DomainEvent)
	}
	o.events[event.EventID()] = event
}

// RemoveDomainEvent removes a previously registered domain event by its EventID.
func (o *AggregateRoot) RemoveDomainEvent(event DomainEvent) {
	delete(o.events, event.EventID())
}

// ClearDomainEvent discards all pending domain events, typically called after events
// have been dispatched.
func (o *AggregateRoot) ClearDomainEvent() {
	o.events = make(map[string]DomainEvent)
}
