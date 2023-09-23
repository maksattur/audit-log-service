package handler

import (
	"encoding/json"
	"github.com/maksattur/audit-log-service/internal"
	"github.com/maksattur/audit-log-service/internal/domain"
	"log"
	"net/http"
	"strconv"
)

func (hh HttpHandler) events(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	filter, err := hh.validateEventsRequest(r)
	if err != nil {
		log.Printf("Requested URL %s, error: %s\n", r.URL.String(), err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	events, err := hh.service.EventsHttp(r.Context(), filter)
	if err != nil {
		log.Printf("Requested URL %s, error %s\n", r.URL.String(), err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := make([]Event, len(events))

	for i, event := range events {
		specific, _ := json.Marshal(event.Specific())
		response[i] = Event{
			Common: CommonFields{
				UserID:    event.UserID(),
				EventType: event.EventType(),
				Timestamp: event.Timestamp(),
			},
			Specific: specific,
		}
	}

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(response)
}

func (hh HttpHandler) validateEventsRequest(r *http.Request) (*domain.Filter, error) {
	userID := r.URL.Query().Get("user_id")
	eventType := r.URL.Query().Get("event_type")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	limitStr := r.URL.Query().Get("limit")

	var limit uint64

	if limitStr != "" {
		var err error
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			return nil, &internal.CustomError{
				OriginalError: err,
				Message:       "error parse limit",
			}
		}
	}

	filter, err := domain.NewFilter(userID, eventType, from, to, limit)
	if err != nil {
		return nil, &internal.CustomError{
			OriginalError: err,
			Message:       "error build filter",
		}
	}

	return filter, nil
}
