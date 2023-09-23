package handler

import (
	"encoding/json"
	"time"
)

type Credentials struct {
	Login    string `json:"Login"`
	Password string `json:"password"`
}

type Event struct {
	Common   CommonFields    `json:"common"`
	Specific json.RawMessage `json:"specific"`
}

type CommonFields struct {
	UserID    string    `json:"user_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
}
