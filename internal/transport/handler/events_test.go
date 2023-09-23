package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/maksattur/audit-log-service/internal"
	"github.com/maksattur/audit-log-service/internal/domain"
	"github.com/maksattur/audit-log-service/internal/transport/handler/mock_handler"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestValidateEventsRequest(t *testing.T) {
	fullFilter, _ := domain.NewFilter("123", "bill", "2023-09-23T16:38:25Z", "2023-09-24T16:38:25Z", 10)
	filterWithoutUserID, _ := domain.NewFilter("", "bill", "2023-09-23T16:38:25Z", "2023-09-24T16:38:25Z", 10)
	filterWithoutUserIDEventType, _ := domain.NewFilter("", "", "2023-09-23T16:38:25Z", "2023-09-24T16:38:25Z", 10)

	testCases := []struct {
		name           string
		queryParams    string
		expectedFilter *domain.Filter
		expectedError  error
	}{
		{
			name:           "Valid Request",
			queryParams:    "user_id=123&event_type=bill&from=2023-09-23T16:38:25Z&to=2023-09-24T16:38:25Z&limit=10",
			expectedFilter: fullFilter,
			expectedError:  nil,
		},
		{
			name:           "Valid Request Without User_ID",
			queryParams:    "event_type=bill&from=2023-09-23T16:38:25Z&to=2023-09-24T16:38:25Z&limit=10",
			expectedFilter: filterWithoutUserID,
			expectedError:  nil,
		},
		{
			name:           "Valid Request Without User_ID",
			queryParams:    "from=2023-09-23T16:38:25Z&to=2023-09-24T16:38:25Z&limit=10",
			expectedFilter: filterWithoutUserIDEventType,
			expectedError:  nil,
		},
		{
			name:           "Filter To Less From",
			queryParams:    "user_id=123&event_type=bill&from=2023-09-24T16:38:25Z&to=2023-09-23T16:38:25Z&limit=10",
			expectedFilter: nil,
			expectedError:  &internal.CustomError{},
		},
		{
			name:           "Missing Params",
			queryParams:    "",
			expectedFilter: nil,
			expectedError:  &internal.CustomError{},
		},
		{
			name:           "Invalid Limit",
			queryParams:    "user_id=123&event_type=login&from=2023-01-01&to=2023-01-31&limit=invalid",
			expectedFilter: nil,
			expectedError:  &internal.CustomError{},
		},
	}

	handler := HttpHandler{}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/events?"+tc.queryParams, nil)
			require.NoError(t, err)

			filter, err := handler.validateEventsRequest(req)

			if tc.expectedFilter != nil {
				require.Equal(t, *tc.expectedFilter, *filter)
				require.NoError(t, err, tc.expectedError)
			} else {
				require.Nil(t, tc.expectedFilter, nil)
				require.Error(t, err)
				require.IsType(t, tc.expectedError, err)
			}

		})
	}
}

func TestEvents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuditLogService := mock_handler.NewMockAuditLogService(ctrl)

	handler := HttpHandler{
		service: mockAuditLogService,
	}
	fromStr := "2023-09-23T16:38:25Z"
	toStr := "2023-09-24T16:38:25Z"

	from, _ := time.Parse(time.RFC3339, fromStr)

	specific := struct {
		Name string
		Age  int
	}{
		Name: "Van der Gam",
		Age:  41,
	}

	event, _ := domain.NewEvent("123", "bill", from, specific)

	testCases := []struct {
		name            string
		queryParams     string
		expectedEvents  []*domain.Event
		expectedStatus  int
		expectedContent string
	}{
		{
			name:            "Valid Request",
			queryParams:     fmt.Sprintf("user_id=123&event_type=bill&from=%s&to=%s&limit=10", fromStr, toStr),
			expectedEvents:  []*domain.Event{event},
			expectedStatus:  http.StatusOK,
			expectedContent: `[{"common":{"user_id":"123","event_type":"bill","timestamp":"2023-09-23T16:38:25Z"},"specific":{"Name":"Van der Gam","Age":41}}]`,
		},
		{
			name:            "Invalid Request",
			queryParams:     "",
			expectedEvents:  nil,
			expectedStatus:  http.StatusBadRequest,
			expectedContent: "",
		},
		{
			name:            "Invalid Limit",
			queryParams:     fmt.Sprintf("user_id=123&event_type=bill&from=%s&to=%s&limit=invalid", fromStr, toStr),
			expectedEvents:  nil,
			expectedStatus:  http.StatusBadRequest,
			expectedContent: "",
		},
		{
			name:            "Internal Error",
			queryParams:     fmt.Sprintf("user_id=123&event_type=bill&from=%s&to=%s&limit=10", fromStr, toStr),
			expectedEvents:  nil,
			expectedStatus:  http.StatusInternalServerError,
			expectedContent: "",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/events?"+tc.queryParams, nil)
			require.NoError(t, err)
			w := httptest.NewRecorder()

			if tc.expectedStatus == http.StatusOK {
				mockAuditLogService.EXPECT().EventsHttp(gomock.Any(), gomock.Any()).Return(tc.expectedEvents, nil).Times(1)

			} else if tc.expectedStatus == http.StatusBadRequest {
				mockAuditLogService.EXPECT().EventsHttp(gomock.Any(), gomock.Any()).Return(tc.expectedEvents, errors.New("bad request")).Times(0)
			} else {
				mockAuditLogService.EXPECT().EventsHttp(gomock.Any(), gomock.Any()).Return(tc.expectedEvents, errors.New("internal error")).Times(1)
			}

			handler.events(w, req)

			if tc.expectedContent != "" {
				require.Equal(t, tc.expectedContent, strings.TrimSpace(w.Body.String()))
				var response []Event
				err = json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)
			}

			require.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
