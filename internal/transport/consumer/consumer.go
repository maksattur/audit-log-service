package consumer

import (
	"context"
	"github.com/maksattur/audit-log-service/internal/domain"
	"log"
)

type ConsumeAdapter interface {
	Receive(ctx context.Context, eventChan chan<- Event, errChan chan<- error)
	Close() error
}

// AuditLogService is audit log service for receive data
type AuditLogService interface {
	SaveData(ctx context.Context, event *domain.Event) error
}

type Consumer struct {
	service  AuditLogService
	consumer ConsumeAdapter
}

func NewConsumer(service AuditLogService, consumer ConsumeAdapter) *Consumer {
	return &Consumer{
		service:  service,
		consumer: consumer,
	}
}

func (c *Consumer) Receive(ctx context.Context) {
	eventChan := make(chan Event)
	errorChan := make(chan error)
	go c.consumer.Receive(ctx, eventChan, errorChan)

	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-eventChan:
			if !ok {
				return
			}
			e, err := eventBrokerToDomain(event)
			if err != nil {
				log.Printf("Error cast: %v\n", err)
				continue
			}
			if err := c.service.SaveData(ctx, e); err != nil {
				log.Printf("Error save to storage: %v\n", err)
			}
		case err, ok := <-errorChan:
			if !ok {
				continue
			}
			log.Printf("Error: %v\n", err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}

func eventBrokerToDomain(e Event) (*domain.Event, error) {
	event, err := domain.NewEvent(e.Common.UserID, e.Common.EventType, e.Common.Timestamp, e.Specific)
	if err != nil {
		return nil, err
	}
	return event, nil
}
