package internalhttp

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/models"
	resp "github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/lib/api/response"
	"golang.org/x/exp/slog"
)

type CreateRequest struct {
	Title            string         `json:"title"`
	Description      *string        `json:"description"`
	UserID           int64          `json:"user_id"`
	StartDate        time.Time      `json:"start_date"`
	EndDate          time.Time      `json:"end_date"`
	NotificationTime *time.Duration `json:"notification_time"`
}

func (h *Handler) createEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.log.With(slog.String("request_id", middleware.GetReqID(r.Context())))

		var event CreateRequest
		if err := parseBody(r, &event); err != nil {
			log.Error("Parse request body", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		if err := event.validate(); err != nil {
			log.Error("Validate event", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		if err := h.app.CreateEvent(r.Context(), event.toModel()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (r *CreateRequest) validate() error {
	if len(r.Title) == 0 {
		return errors.New("field title is empty")
	}
	if r.UserID == 0 {
		return errors.New("field user_id is empty")
	}
	if r.StartDate.IsZero() {
		return errors.New("field start_date is empty")
	}
	if r.EndDate.IsZero() {
		return errors.New("field end_date is empty")
	}
	return nil
}

func (r *CreateRequest) toModel() *models.Event {
	return &models.Event{
		Title:            r.Title,
		Description:      r.Description,
		UserID:           r.UserID,
		StartDate:        r.StartDate,
		EndDate:          r.EndDate,
		NotificationTime: r.NotificationTime,
	}
}

type Event struct {
	ID               string         `json:"id"`
	Title            string         `json:"title"`
	Description      *string        `json:"description"`
	UserID           int64          `json:"user_id"`
	StartDate        time.Time      `json:"start_date"`
	EndDate          time.Time      `json:"end_date"`
	NotificationTime *time.Duration `json:"notification_time"`
}

func (h *Handler) updateEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.log.With(slog.String("request_id", middleware.GetReqID(r.Context())))

		eventID := parseID(r)

		var event Event
		if err := parseBody(r, &event); err != nil {
			log.Error("Parse request body", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		event.ID = eventID
		if err := h.app.UpdateEvent(r.Context(), event.toModel()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (r *Event) toModel() *models.Event {
	return &models.Event{
		ID:               r.ID,
		Title:            r.Title,
		Description:      r.Description,
		UserID:           r.UserID,
		StartDate:        r.StartDate,
		EndDate:          r.EndDate,
		NotificationTime: r.NotificationTime,
	}
}

type StartDate time.Time

type GetByDateRequest struct {
	UserID int64     `json:"user_id"`
	Date   StartDate `json:"start_date"`
}

type EventsResponse []Event

func (h *Handler) getEventsByDay() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.log.With(slog.String("request_id", middleware.GetReqID(r.Context())))

		var request GetByDateRequest
		if err := parseBody(r, &request); err != nil {
			log.Error("Parse request body", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		if err := request.validate(); err != nil {
			log.Error("Validate get by day request", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		events, err := h.app.GetEventByDay(r.Context(), request.UserID, time.Time(request.Date))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, toResponse(events))
	}
}

func (h *Handler) getEventsByWeek() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.log.With(slog.String("request_id", middleware.GetReqID(r.Context())))

		var request GetByDateRequest
		if err := parseBody(r, &request); err != nil {
			log.Error("Parse request body", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		if err := request.validate(); err != nil {
			log.Error("Validate get by week request", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		events, err := h.app.GetEventByWeek(r.Context(), request.UserID, time.Time(request.Date))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, toResponse(events))
	}
}

func (h *Handler) getEventsByMonth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.log.With(slog.String("request_id", middleware.GetReqID(r.Context())))

		var request GetByDateRequest
		if err := parseBody(r, &request); err != nil {
			log.Error("Parse request body", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		if err := request.validate(); err != nil {
			log.Error("Validate get by month request", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		events, err := h.app.GetEventByMonth(r.Context(), request.UserID, time.Time(request.Date))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, toResponse(events))
	}
}

func toResponse(events []models.Event) EventsResponse {
	resp := make(EventsResponse, len(events))
	for i := range events {
		resp[i] = Event{
			ID:               events[i].ID,
			Title:            events[i].Title,
			Description:      events[i].Description,
			UserID:           events[i].UserID,
			StartDate:        events[i].StartDate,
			EndDate:          events[i].EndDate,
			NotificationTime: events[i].NotificationTime,
		}
	}
	return resp
}

func (r *GetByDateRequest) validate() error {
	if r.UserID == 0 {
		return errors.New("field user_id is empty")
	}
	if time.Time(r.Date).IsZero() {
		return errors.New("field start_date is empty")
	}
	return nil
}

func (sd *StartDate) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`)
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return err
	}
	*sd = StartDate(t)
	return nil
}

func (h *Handler) deleteEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventID := parseID(r)

		if err := h.app.DeleteEvent(r.Context(), eventID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
