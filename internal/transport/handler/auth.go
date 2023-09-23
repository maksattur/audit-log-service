package handler

import (
	"encoding/json"
	"errors"
	"github.com/maksattur/audit-log-service/internal"
	"log"
	"net/http"
)

const (
	Login        = "bambino"
	Password     = "qwerty123"
	PasswordHash = "$2a$10$wRI1uRkdb3uxw9dJyHrudeSGFPPo5aIFO4LanU.GAq0YfknryquFW"
)

func (hh HttpHandler) auth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	credentials := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		log.Printf("Requested URL %s, error: %s\n", r.URL.String(), err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := hh.validateAuthRequest(credentials, Login, PasswordHash); err != nil {
		log.Printf("Requested URL %s, error: %s\n", r.URL.String(), err.Error())
		switch {
		case errors.Is(err, ErrLoginOrPasswordIsEmpty):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, ErrLoginOrPasswordIncorrect):
			w.WriteHeader(http.StatusUnauthorized)
		}
		return
	}

	token, err := hh.token.BuildToken()
	if err != nil {
		err = &internal.CustomError{
			OriginalError: err,
			Message:       "build token",
		}
		log.Printf("Requested URL %s, error: %s\n", r.URL.String(), err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (hh HttpHandler) validateAuthRequest(cred *Credentials, login string, passHash string) error {
	if cred.Login == "" || cred.Password == "" {
		return ErrLoginOrPasswordIsEmpty
	}
	if cred.Login != login || !CheckPasswordHash(cred.Password, passHash) {
		return ErrLoginOrPasswordIncorrect
	}
	return nil
}
