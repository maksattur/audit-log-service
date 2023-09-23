package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/maksattur/audit-log-service/internal/transport/handler/mock_handler"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateAuthRequest(t *testing.T) {
	handler := HttpHandler{}

	cred := &Credentials{
		Login:    "alex",
		Password: "123",
	}
	passHash, _ := GeneratePasswordHash(cred.Password)
	err := handler.validateAuthRequest(cred, cred.Login, passHash)
	require.NoError(t, err)
}

func TestValidateAuthRequestError(t *testing.T) {
	testCases := []struct {
		name     string
		cred     *Credentials
		login    string
		passHash string
		expErr   error
	}{
		{
			name: "Login Is Empty",
			cred: &Credentials{
				Login: "",
			},
			expErr: ErrLoginOrPasswordIsEmpty,
		},
		{
			name: "Password Is Empty",
			cred: &Credentials{
				Password: "",
			},
			expErr: ErrLoginOrPasswordIsEmpty,
		},
		{
			name: "Login Is Incorrect",
			cred: &Credentials{
				Login:    "alex",
				Password: "qwerty123",
			},
			login:  "sam",
			expErr: ErrLoginOrPasswordIncorrect,
		},
		{
			name: "Password Is Incorrect",
			cred: &Credentials{
				Login:    "alex",
				Password: "123123",
			},
			passHash: "$2a$10$wRI1uRkdb3uxw9dJyHrudeSGFPPo5aIFO4LanU.GAq0YfknryquFW",
			expErr:   ErrLoginOrPasswordIncorrect,
		},
	}

	handler := HttpHandler{}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			err := handler.validateAuthRequest(tCase.cred, tCase.login, tCase.passHash)
			require.Error(t, err)
			require.ErrorIs(t, err, tCase.expErr)
		})
	}
}

func TestAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenizer := mock_handler.NewMockTokenizer(ctrl)

	handler := HttpHandler{
		token: mockTokenizer,
	}

	testCases := []struct {
		name            string
		requestBody     interface{}
		buildTokenError error
		expectedStatus  int
		expectedBody    string
	}{
		{
			name:            "Valid Authentication",
			requestBody:     &Credentials{Login: Login, Password: Password},
			buildTokenError: nil,
			expectedStatus:  http.StatusOK,
			expectedBody:    `{"token":"test_token"}`,
		},
		{
			name:            "Login Is Empty",
			requestBody:     &Credentials{Login: "", Password: Password},
			buildTokenError: nil,
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    "",
		},
		{
			name:            "Password Is Empty",
			requestBody:     &Credentials{Login: Login, Password: ""},
			buildTokenError: nil,
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    "",
		},
		{
			name:            "Login Is Incorrect",
			requestBody:     &Credentials{Login: "123", Password: Password},
			buildTokenError: nil,
			expectedStatus:  http.StatusUnauthorized,
			expectedBody:    "",
		},
		{
			name:            "Password Is Incorrect",
			requestBody:     &Credentials{Login: Login, Password: "123"},
			buildTokenError: nil,
			expectedStatus:  http.StatusUnauthorized,
			expectedBody:    "",
		},
		{
			name:            "Invalid JSON Request Body",
			requestBody:     "str",
			buildTokenError: nil,
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    "",
		},
		{
			name:            "Token Build Error",
			requestBody:     &Credentials{Login: Login, Password: Password},
			buildTokenError: errors.New("build error"),
			expectedStatus:  http.StatusInternalServerError,
			expectedBody:    "",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest("POST", "/auth", bytes.NewReader(reqBody))
			w := httptest.NewRecorder()

			if tc.expectedStatus == http.StatusOK {
				mockTokenizer.EXPECT().BuildToken().Return("test_token", tc.buildTokenError).Times(1)
			} else if tc.expectedStatus == http.StatusInternalServerError {
				mockTokenizer.EXPECT().BuildToken().Return("", tc.buildTokenError).Times(1)
			} else {
				mockTokenizer.EXPECT().BuildToken().Return("", nil).Times(0)
			}

			handler.auth(w, req)

			require.Equal(t, tc.expectedStatus, w.Code)
			require.Equal(t, tc.expectedBody, strings.TrimSpace(w.Body.String()))
		})
	}
}
