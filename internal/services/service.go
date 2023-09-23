package services

import (
	"context"
	"github.com/maksattur/audit-log-service/internal/domain"
)

// EventRepositoryAdapter is an event repository adapter for storages
type EventRepositoryAdapter interface {
	SaveEvent(ctx context.Context, event *domain.Event) error
	Events(ctx context.Context, filter *domain.Filter) ([]*domain.Event, error)
}

type EventService struct {
	repo EventRepositoryAdapter
}

func NewEventService(repo EventRepositoryAdapter) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (es *EventService) EventsHttp(ctx context.Context, filter *domain.Filter) ([]*domain.Event, error) {
	return es.repo.Events(ctx, filter)
}

func (es *EventService) SaveData(ctx context.Context, event *domain.Event) error {
	return es.repo.SaveEvent(ctx, event)
}
