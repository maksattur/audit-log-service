package handler

import (
	"log"
	"net/http"
	"strings"
)

// middleware is JWT token verifier
func (hh HttpHandler) middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			log.Printf("Requested URL %s, header Authorization is empty\n", r.URL.String())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			log.Printf("Requested URL %s, header Authorization has not prefix 'Bearer '\n", r.URL.String())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authorizationHeader, "Bearer ")
		if err := hh.token.VerifyToken(token); err != nil {
			log.Printf("Requested URL %s, error %s\n", r.URL.String(), err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}
