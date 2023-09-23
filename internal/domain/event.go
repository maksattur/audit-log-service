package domain

import (
	"fmt"
	"time"
)

type Event struct {
	userID    string
	eventType string
	timestamp time.Time
	specific  any
}

func NewEvent(userID string, eventType string, timestamp time.Time, specific any) (*Event, error) {
	if userID == "" {
		return nil, fmt.Errorf("%w: user id is required", ErrRequired)
	}
	if eventType == "" {
		return nil, fmt.Errorf("%w: event type is required", ErrRequired)
	}
	if timestamp.IsZero() {
		return nil, fmt.Errorf("%w: event timestamp is required", ErrRequired)
	}
	return &Event{
		userID:    userID,
		eventType: eventType,
		timestamp: timestamp,
		specific:  specific,
	}, nil
}

func (e *Event) UserID() string {
	return e.userID
}

func (e *Event) EventType() string {
	return e.eventType
}

func (e *Event) Timestamp() time.Time {
	return e.timestamp
}

func (e *Event) Specific() any {
	return e.specific
}
