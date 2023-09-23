package ch

import (
	"encoding/json"
	"github.com/maksattur/audit-log-service/internal/domain"
	"time"
)

type Event struct {
	UserID    string    `db:"user_id"`
	EventType string    `db:"event_type"`
	Timestamp time.Time `db:"timestamp"`
	Specific  string    `db:"specific"`
	CreatedAt time.Time `db:"created_at"`
}

func eventDomainToStore(e *domain.Event) (*Event, error) {
	specificBytes, err := json.Marshal(e.Specific())
	if err != nil {
		return nil, err
	}
	return &Event{
		UserID:    e.UserID(),
		EventType: e.EventType(),
		Timestamp: e.Timestamp(),
		Specific:  string(specificBytes),
	}, nil
}
