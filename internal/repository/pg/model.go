package pg

import (
	"github.com/maksattur/audit-log-service/internal/domain"
	"time"
)

type Event struct {
	ID        int       `db:"id"`
	UserID    string    `db:"user_id"`
	EventType string    `db:"event_type"`
	Timestamp time.Time `db:"timestamp"`
	Specific  any       `db:"specific"`
	CreatedAt time.Time `db:"created_at"`
}

func eventDomainToStore(e *domain.Event) *Event {
	return &Event{
		UserID:    e.UserID(),
		EventType: e.EventType(),
		Timestamp: e.Timestamp(),
		Specific:  e.Specific(),
	}
}
