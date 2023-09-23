package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	PingRoute = "/ping"
)

// HttpHandler is HTTP handler for audit log service
type HttpHandler struct {
}

func NewHttpHandler() *HttpHandler {
	return &HttpHandler{}
}

func (hh HttpHandler) Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc(PingRoute, hh.ping).Methods(http.MethodGet)

	return router
}

func (hh HttpHandler) ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application-json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]bool{
		"pong": true,
	})
}
