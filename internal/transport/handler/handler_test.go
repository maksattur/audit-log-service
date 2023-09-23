package handler

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/maksattur/audit-log-service/internal/transport/handler/mock_handler"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHttpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenizer := mock_handler.NewMockTokenizer(ctrl)
	mockAuditLogService := mock_handler.NewMockAuditLogService(ctrl)

	handler := NewHttpHandler(mockAuditLogService, mockTokenizer)

	req, err := http.NewRequest("POST", AuthRoute, bytes.NewReader([]byte("str")))
	require.NoError(t, err)

	router := mux.NewRouter()
	router.HandleFunc(AuthRoute, handler.auth).Methods(http.MethodPost)
	w := httptest.NewRecorder()

	mockTokenizer.EXPECT().BuildToken().Return("", nil).Times(0)
	mockAuditLogService.EXPECT().EventsHttp(gomock.Any(), gomock.Any()).Return(nil, errors.New("bad request")).Times(0)

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenizer := mock_handler.NewMockTokenizer(ctrl)
	mockAuditLogService := mock_handler.NewMockAuditLogService(ctrl)

	handler := NewHttpHandler(mockAuditLogService, mockTokenizer)

	ts := httptest.NewServer(handler.Router())
	defer ts.Close()

	resp, err := http.Post(ts.URL+AuthRoute, "application/json", nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
