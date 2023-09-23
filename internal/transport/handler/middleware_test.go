package handler

import (
	"github.com/golang/mock/gomock"
	"github.com/maksattur/audit-log-service/internal/token_manager"
	"github.com/maksattur/audit-log-service/internal/transport/handler/mock_handler"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mockTokenizer := mock_handler.NewMockTokenizer(ctrl)

	hh := HttpHandler{
		token: mockTokenizer,
	}

	server := httptest.NewServer(hh.middleware(handler))
	defer server.Close()

	t.Run("No Authorization Header", func(t *testing.T) {
		req, err := http.NewRequest("GET", server.URL, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Invalid Authorization Header", func(t *testing.T) {
		req, err := http.NewRequest("GET", server.URL, nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "InvalidToken")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Token Is Expired", func(t *testing.T) {
		mockTokenizer.EXPECT().VerifyToken(gomock.Any()).Return(token_manager.ErrTokenIsExpired)

		req, err := http.NewRequest("GET", server.URL, nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer ValidToken")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Valid Authorization Header", func(t *testing.T) {
		mockTokenizer.EXPECT().VerifyToken(gomock.Any()).Return(nil)

		req, err := http.NewRequest("GET", server.URL, nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer ValidToken")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
