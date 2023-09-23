package ch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/maksattur/audit-log-service/internal/config"
	"github.com/maksattur/audit-log-service/internal/domain"
)

type ClickHouse struct {
	conn driver.Conn
}

func NewClickHouse(ctx context.Context, cfg config.ClickHouse) (*ClickHouse, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.ClickHouseAddr},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.Username,
			Password: cfg.Password,
		},
	})
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}

	return &ClickHouse{
		conn: conn,
	}, nil
}

// SaveEvent is to save event to ClickHouse
func (c *ClickHouse) SaveEvent(ctx context.Context, event *domain.Event) error {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)

	e, err := eventDomainToStore(event)
	if err != nil {
		return fmt.Errorf("marshal specific: %w", err)
	}

	sql, args, err := queryBuilder.
		Insert("event").Columns("user_id", "event_type", "timestamp", "specific").
		Values(e.UserID, e.EventType, e.Timestamp, e.Specific).ToSql()
	if err != nil {
		return fmt.Errorf("build sql: %w", err)
	}

	err = c.conn.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("insert data: %w", err)
	}

	return nil
}

// Events get events from ClickHouse
func (c *ClickHouse) Events(ctx context.Context, filter *domain.Filter) ([]*domain.Event, error) {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	query := queryBuilder.
		Select("user_id", "event_type", "timestamp", "specific").
		From("event").
		Where("timestamp >= ?", filter.From()).
		Where("timestamp <= ?", filter.To()).Limit(filter.Limit())

	if filter.UserID() != "" {
		query = query.Where(squirrel.Eq{"user_id": filter.UserID()})
	}

	if filter.EventType() != "" {
		query = query.Where(squirrel.Eq{"event_type": filter.EventType()})
	}

	sql, args, _ := query.ToSql()

	rows, err := c.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("perform sql query: %w", err)
	}
	defer rows.Close()

	var result []*domain.Event

	for rows.Next() {
		cre, err := c.scanEvent(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, cre)
	}

	return result, nil
}

func (c *ClickHouse) scanEvent(row pgx.Row) (*domain.Event, error) {
	var (
		event    Event
		specific any
	)
	if err := row.Scan(&event.UserID, &event.EventType, &event.Timestamp, &event.Specific); err != nil {
		return nil, fmt.Errorf("scan row: %w", err)
	}

	if err := json.Unmarshal([]byte(event.Specific), &specific); err != nil {
		return nil, fmt.Errorf("unmarshal row: %w", err)
	}

	return domain.NewEvent(event.UserID, event.EventType, event.Timestamp, specific)
}
