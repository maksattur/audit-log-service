package consumer

import (
	"encoding/json"
	"time"
)

type Event struct {
	Common   CommonFields    `json:"common"`
	Specific json.RawMessage `json:"specific"`
}

type CommonFields struct {
	UserID    string    `json:"user_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
}
