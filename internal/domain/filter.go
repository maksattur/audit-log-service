package domain

import (
	"fmt"
	"time"
)

const defaultLimit = 10000

type Filter struct {
	userID    string
	eventType string
	from      time.Time
	to        time.Time
	limit     uint64
}

func NewFilter(userID string, eventType string, from string, to string, limit uint64) (*Filter, error) {
	format := time.RFC3339
	parsedFrom, err := time.Parse(format, from)
	if err != nil {
		return nil, fmt.Errorf("%w: 'from' date is required, parse %w", ErrRequired, err)
	}

	parsedTo, err := time.Parse(format, to)
	if err != nil {
		return nil, fmt.Errorf("%w: to date is required, parse %w", ErrRequired, err)
	}

	if parsedTo.Before(parsedFrom) {
		return nil, fmt.Errorf("%w: 'to' less 'from'", ErrDateFormat)
	}

	if limit > defaultLimit || limit == 0 {
		limit = defaultLimit
	}

	return &Filter{
		userID:    userID,
		eventType: eventType,
		from:      parsedFrom,
		to:        parsedTo,
		limit:     limit,
	}, nil
}

func (f *Filter) UserID() string {
	return f.userID
}

func (f *Filter) EventType() string {
	return f.eventType
}

func (f *Filter) From() time.Time {
	return f.from
}

func (f *Filter) To() time.Time {
	return f.to
}

func (f *Filter) Limit() uint64 {
	return f.limit
}
