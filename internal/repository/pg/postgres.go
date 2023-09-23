package pg

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/maksattur/audit-log-service/internal/config"
	"github.com/maksattur/audit-log-service/internal/domain"
	"time"
)

type Postgres struct {
	conn *pgxpool.Pool
}

func NewPostgres(ctx context.Context, cfg config.Postgres) (*Postgres, error) {
	connPool, err := pgxpool.Connect(ctx, cfg.DSN)
	if err != nil {
		return nil, err
	}
	return &Postgres{
		conn: connPool,
	}, nil
}

func (p *Postgres) SaveEvent(ctx context.Context, event *domain.Event) error {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	e := eventDomainToStore(event)
	sql, args, err := queryBuilder.
		Insert("event").Columns("user_id", "event_type", "timestamp", "specific").
		Values(e.UserID, e.EventType, e.Timestamp, e.Specific).ToSql()
	if err != nil {
		return fmt.Errorf("build sql: %w", err)
	}
	_, err = p.conn.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("insert data: %w", err)
	}
	return nil
}

func (p *Postgres) Events(ctx context.Context, filter *domain.Filter) ([]*domain.Event, error) {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query := queryBuilder.
		Select("user_id", "event_type", "timestamp", "specific").
		From("event").
		Where("timestamp >= ?", filter.From()).
		Where("timestamp <= ?", filter.To()).Limit(10000)

	if filter.UserID() != "" {
		query = query.Where(squirrel.Eq{"user_id": filter.UserID()})
	}

	if filter.EventType() != "" {
		query = query.Where(squirrel.Eq{"event_type": filter.EventType()})
	}

	sql, args, _ := query.ToSql()

	rows, err := p.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("perform sql query: %w", err)
	}
	defer rows.Close()

	var result []*domain.Event

	for rows.Next() {
		cre, err := p.scanEvent(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, cre)
	}

	return result, nil
}

func (p *Postgres) scanEvent(row pgx.Row) (*domain.Event, error) {
	type Event struct {
		UserID    string
		EventType string
		Timestamp time.Time
		Specific  any
	}

	var e Event

	if err := row.Scan(&e.UserID, &e.EventType, &e.Timestamp, &e.Specific); err != nil {
		return nil, fmt.Errorf("scan row: %w", err)
	}

	return domain.NewEvent(e.UserID, e.EventType, e.Timestamp, e.Specific)
}
