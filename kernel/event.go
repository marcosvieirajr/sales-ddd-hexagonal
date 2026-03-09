package kernel

import "time"

type Event struct {
	DateOccurred time.Time `json:"occurred_at"`
}

func (e Event) OccurredAt() time.Time {
	return e.DateOccurred
}
