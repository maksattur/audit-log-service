package handler

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/maksattur/audit-log-service/internal/domain"
	"net/http"
)

const (
	AuthRoute  = "/auth"
	EventRoute = "/events"
)

// AuditLogService is audit log service
type AuditLogService interface {
	EventsHttp(ctx context.Context, filter *domain.Filter) ([]*domain.Event, error)
}

// Tokenizer for token manager
type Tokenizer interface {
	BuildToken() (string, error)
	VerifyToken(string) error
}

// HttpHandler is HTTP handler for audit log service
type HttpHandler struct {
	service AuditLogService
	token   Tokenizer
}

func NewHttpHandler(service AuditLogService, token Tokenizer) *HttpHandler {
	return &HttpHandler{
		service: service,
		token:   token,
	}
}

func (hh HttpHandler) Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc(AuthRoute, hh.auth).Methods(http.MethodPost)

	eventsRouter := router.PathPrefix(EventRoute).Subrouter()
	eventsRouter.Use(hh.middleware)
	eventsRouter.HandleFunc("", hh.events).Methods(http.MethodGet)

	return router
}
