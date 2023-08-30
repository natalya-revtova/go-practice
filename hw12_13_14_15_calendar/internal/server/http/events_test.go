package internalhttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/logger"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/models"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/server/mocks"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateHandler(t *testing.T) {
	cases := []struct {
		name      string
		event     CreateRequest
		body      map[string]interface{}
		code      int
		respError string
		mockError error
	}{
		{
			name: "success",
			event: CreateRequest{
				Title:            "test",
				Description:      nil,
				UserID:           1,
				StartDate:        time.Date(2023, 8, 16, 12, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2023, 8, 16, 13, 0, 0, 0, time.UTC),
				NotificationTime: nil,
			},
			body: map[string]interface{}{
				"title":      "test",
				"user_id":    1,
				"start_date": "2023-08-16T12:00:00Z",
				"end_date":   "2023-08-16T13:00:00Z",
			},
			code: http.StatusCreated,
		},
		{
			name:      "empty body",
			body:      nil,
			respError: "request body is empty",
			code:      http.StatusBadRequest,
		},
		{
			name: "empty title",
			body: map[string]interface{}{
				"user_id":    1,
				"start_date": "2023-08-16T12:00:00Z",
				"end_date":   "2023-08-16T13:00:00Z",
			},
			respError: "field title is empty",
			code:      http.StatusBadRequest,
		},
		{
			name: "empty user_id",
			body: map[string]interface{}{
				"title":      "test",
				"start_date": "2023-08-16T12:00:00Z",
				"end_date":   "2023-08-16T13:00:00Z",
			},
			respError: "field user_id is empty",
			code:      http.StatusBadRequest,
		},
		{
			name: "empty start_date",
			body: map[string]interface{}{
				"title":    "test",
				"user_id":  1,
				"end_date": "2023-08-16T13:00:00Z",
			},
			respError: "field start_date is empty",
			code:      http.StatusBadRequest,
		},
		{
			name: "empty end_date",
			body: map[string]interface{}{
				"title":      "test",
				"user_id":    1,
				"start_date": "2023-08-16T12:00:00Z",
			},
			respError: "field end_date is empty",
			code:      http.StatusBadRequest,
		},
		{
			name: "create event error",
			event: CreateRequest{
				Title:            "test",
				Description:      nil,
				UserID:           1,
				StartDate:        time.Date(2023, 8, 16, 12, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2023, 8, 16, 13, 0, 0, 0, time.UTC),
				NotificationTime: nil,
			},
			body: map[string]interface{}{
				"title":      "test",
				"user_id":    1,
				"start_date": "2023-08-16T12:00:00Z",
				"end_date":   "2023-08-16T13:00:00Z",
			},
			mockError: errors.New("unexpected error"),
			code:      http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			appMock := mocks.NewCalendar(t)

			if tc.respError == "" || tc.mockError != nil {
				appMock.On("CreateEvent", mock.Anything, tc.event.toModel()).
					Return(tc.mockError).
					Once()
			}

			handler := chi.NewRouter()
			handler.Post(eventsURL, NewHandler(logger.NewMock(), appMock).createEvent())

			var body []byte
			var err error
			if tc.body != nil {
				body, err = json.Marshal(tc.body)
				require.NoError(t, err)
			}

			req, err := http.NewRequest(http.MethodPost, eventsURL, bytes.NewReader(body))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.code, rr.Code)
		})
	}
}

func TestUpdateHandler(t *testing.T) {
	cases := []struct {
		name      string
		event     Event
		body      map[string]interface{}
		code      int
		respError string
		mockError error
	}{
		{
			name: "success",
			event: Event{
				ID:    "id-1",
				Title: "test",
			},
			body: map[string]interface{}{
				"id":    "id-1",
				"title": "test",
			},
			code: http.StatusOK,
		},
		{
			name: "empty body",
			event: Event{
				ID:    "id-1",
				Title: "test",
			},
			body:      nil,
			respError: "request body is empty",
			code:      http.StatusBadRequest,
		},
		{
			name: "update event error",
			event: Event{
				ID:    "id-1",
				Title: "test",
			},
			body: map[string]interface{}{
				"id":    "id-1",
				"title": "test",
			},
			mockError: errors.New("unexpected error"),
			code:      http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			appMock := mocks.NewCalendar(t)

			if tc.respError == "" || tc.mockError != nil {
				appMock.On("UpdateEvent", mock.Anything, tc.event.toModel()).
					Return(tc.mockError).
					Once()
			}

			handler := chi.NewRouter()
			handler.Patch(eventsURL+"/{id}", NewHandler(logger.NewMock(), appMock).updateEvent())

			var body []byte
			var err error
			if tc.body != nil {
				body, err = json.Marshal(tc.body)
				require.NoError(t, err)
			}

			req, err := http.NewRequest(http.MethodPatch, eventsURL+"/"+tc.event.ID, bytes.NewReader(body))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.code, rr.Code)
		})
	}
}

func TestGetByDayHandler(t *testing.T) {
	cases := []struct {
		name      string
		userID    int64
		day       time.Time
		request   GetByDateRequest
		body      map[string]interface{}
		events    []models.Event
		code      int
		respError string
		mockError error
	}{
		{
			name:   "success",
			userID: 1,
			day:    time.Date(2023, 8, 16, 0, 0, 0, 0, time.UTC),
			body: map[string]interface{}{
				"user_id":    1,
				"start_date": "2023-08-16",
			},
			events: []models.Event{
				{
					Title:            "test",
					Description:      nil,
					UserID:           1,
					StartDate:        time.Date(2023, 8, 16, 12, 0, 0, 0, time.UTC),
					EndDate:          time.Date(2023, 8, 16, 13, 0, 0, 0, time.UTC),
					NotificationTime: nil,
					Day:              time.Date(2023, 8, 16, 0, 0, 0, 0, time.UTC),
					Week:             time.Date(2023, 8, 14, 0, 0, 0, 0, time.UTC),
					Month:            time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			code: http.StatusOK,
		},
		{
			name:      "empty body",
			body:      nil,
			respError: "request body is empty",
			code:      http.StatusBadRequest,
		},
		{
			name: "empty user_id",
			body: map[string]interface{}{
				"start_date": "2023-08-16",
			},
			respError: "field user_id is empty",
			code:      http.StatusBadRequest,
		},
		{
			name: "empty start_date",
			body: map[string]interface{}{
				"user_id": 1,
			},
			respError: "json parse error",
			code:      http.StatusBadRequest,
		},
		{
			name: "invalid start_date",
			body: map[string]interface{}{
				"user_id":    1,
				"start_date": "11111",
			},
			respError: "json parse error",
			code:      http.StatusBadRequest,
		},
		{
			name:   "get events error",
			userID: 1,
			day:    time.Date(2023, 8, 16, 0, 0, 0, 0, time.UTC),
			body: map[string]interface{}{
				"user_id":    1,
				"start_date": "2023-08-16",
			},
			mockError: errors.New("unexpected error"),
			code:      http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			appMock := mocks.NewCalendar(t)

			if tc.respError == "" || tc.mockError != nil {
				appMock.On("GetEventByDay", mock.Anything, tc.userID, tc.day).
					Return(tc.events, tc.mockError).
					Once()
			}

			handler := chi.NewRouter()
			handler.Get(eventsURL+"/day", NewHandler(logger.NewMock(), appMock).getEventsByDay())

			var requestBody []byte
			var err error
			if tc.body != nil {
				requestBody, err = json.Marshal(tc.body)
				require.NoError(t, err)
			}

			req, err := http.NewRequest(http.MethodGet, eventsURL+"/day", bytes.NewReader(requestBody))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.code, rr.Code)

			if tc.events != nil {
				var responseBody EventsResponse
				err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				require.Equal(t, toResponse(tc.events), responseBody)
			}
		})
	}
}

func TestGetByWeekHandler(t *testing.T) {
	cases := []struct {
		name      string
		userID    int64
		day       time.Time
		request   GetByDateRequest
		body      map[string]interface{}
		events    []models.Event
		code      int
		respError string
		mockError error
	}{
		{
			name:   "success",
			userID: 1,
			day:    time.Date(2023, 8, 14, 0, 0, 0, 0, time.UTC),
			body: map[string]interface{}{
				"user_id":    1,
				"start_date": "2023-08-14",
			},
			events: []models.Event{
				{
					Title:            "test",
					Description:      nil,
					UserID:           1,
					StartDate:        time.Date(2023, 8, 16, 12, 0, 0, 0, time.UTC),
					EndDate:          time.Date(2023, 8, 16, 13, 0, 0, 0, time.UTC),
					NotificationTime: nil,
					Day:              time.Date(2023, 8, 16, 0, 0, 0, 0, time.UTC),
					Week:             time.Date(2023, 8, 14, 0, 0, 0, 0, time.UTC),
					Month:            time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			code: http.StatusOK,
		},
		{
			name:      "empty body",
			body:      nil,
			respError: "request body is empty",
			code:      http.StatusBadRequest,
		},
		{
			name: "empty user_id",
			body: map[string]interface{}{
				"start_date": "2023-08-14",
			},
			respError: "field user_id is empty",
			code:      http.StatusBadRequest,
		},
		{
			name: "empty start_date",
			body: map[string]interface{}{
				"user_id": 1,
			},
			respError: "field start_date is empty",
			code:      http.StatusBadRequest,
		},
		{
			name: "invalid start_date",
			body: map[string]interface{}{
				"user_id":    1,
				"start_date": "11111",
			},
			respError: "json parse error",
			code:      http.StatusBadRequest,
		},
		{
			name:   "get events error",
			userID: 1,
			day:    time.Date(2023, 8, 14, 0, 0, 0, 0, time.UTC),
			body: map[string]interface{}{
				"user_id":    1,
				"start_date": "2023-08-14",
			},
			mockError: errors.New("unexpected error"),
			code:      http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			appMock := mocks.NewCalendar(t)

			if tc.respError == "" || tc.mockError != nil {
				appMock.On("GetEventByWeek", mock.Anything, tc.userID, tc.day).
					Return(tc.events, tc.mockError).
					Once()
			}

			handler := chi.NewRouter()
			handler.Get(eventsURL+"/"+"week", NewHandler(logger.NewMock(), appMock).getEventsByWeek())

			var body []byte
			var err error
			if tc.body != nil {
				body, err = json.Marshal(tc.body)
				require.NoError(t, err)
			}

			req, err := http.NewRequest(http.MethodGet, eventsURL+"/"+"week", bytes.NewReader(body))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.code, rr.Code)

			if tc.events != nil {
				var responseBody EventsResponse
				err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				require.Equal(t, toResponse(tc.events), responseBody)
			}
		})
	}
}

func TestGetByMonthHandler(t *testing.T) {
	cases := []struct {
		name      string
		userID    int64
		day       time.Time
		request   GetByDateRequest
		body      map[string]interface{}
		events    []models.Event
		code      int
		respError string
		mockError error
	}{
		{
			name:   "success",
			userID: 1,
			day:    time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			body: map[string]interface{}{
				"user_id":    1,
				"start_date": "2023-08-01",
			},
			events: []models.Event{
				{
					Title:            "test",
					Description:      nil,
					UserID:           1,
					StartDate:        time.Date(2023, 8, 16, 12, 0, 0, 0, time.UTC),
					EndDate:          time.Date(2023, 8, 16, 13, 0, 0, 0, time.UTC),
					NotificationTime: nil,
					Day:              time.Date(2023, 8, 16, 0, 0, 0, 0, time.UTC),
					Week:             time.Date(2023, 8, 14, 0, 0, 0, 0, time.UTC),
					Month:            time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			code: http.StatusOK,
		},
		{
			name:      "empty body",
			body:      nil,
			respError: "request body is empty",
			code:      http.StatusBadRequest,
		},
		{
			name: "empty user_id",
			body: map[string]interface{}{
				"start_date": "2023-08-01",
			},
			respError: "field user_id is empty",
			code:      http.StatusBadRequest,
		},
		{
			name: "empty start_date",
			body: map[string]interface{}{
				"user_id": 1,
			},
			respError: "field start_date is empty",
			code:      http.StatusBadRequest,
		},
		{
			name: "invalid start_date",
			body: map[string]interface{}{
				"user_id":    1,
				"start_date": "11111",
			},
			respError: "json parse error",
			code:      http.StatusBadRequest,
		},
		{
			name:   "get events error",
			userID: 1,
			day:    time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			body: map[string]interface{}{
				"user_id":    1,
				"start_date": "2023-08-01",
			},
			mockError: errors.New("unexpected error"),
			code:      http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			appMock := mocks.NewCalendar(t)

			if tc.respError == "" || tc.mockError != nil {
				appMock.On("GetEventByMonth", mock.Anything, tc.userID, tc.day).
					Return(tc.events, tc.mockError).
					Once()
			}

			handler := chi.NewRouter()
			handler.Get(eventsURL+"/"+"month", NewHandler(logger.NewMock(), appMock).getEventsByMonth())

			var body []byte
			var err error
			if tc.body != nil {
				body, err = json.Marshal(tc.body)
				require.NoError(t, err)
			}

			req, err := http.NewRequest(http.MethodGet, eventsURL+"/"+"month", bytes.NewReader(body))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.code, rr.Code)

			if tc.events != nil {
				var responseBody EventsResponse
				err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
				require.NoError(t, err)

				require.Equal(t, toResponse(tc.events), responseBody)
			}
		})
	}
}

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		eventID   string
		code      int
		mockError error
	}{
		{
			name:    "success",
			eventID: "id-1",
			code:    http.StatusNoContent,
		},
		{
			name:      "not existing event",
			eventID:   "id-1",
			mockError: storage.ErrEventNotExist,
			code:      http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			appMock := mocks.NewCalendar(t)

			appMock.On("DeleteEvent", mock.Anything, tc.eventID).
				Return(tc.mockError).
				Once()

			handler := chi.NewRouter()
			handler.Delete(eventsURL+"/{id}", NewHandler(logger.NewMock(), appMock).deleteEvent())

			req, err := http.NewRequest(http.MethodDelete, eventsURL+"/"+tc.eventID, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.code, rr.Code)
		})
	}
}
