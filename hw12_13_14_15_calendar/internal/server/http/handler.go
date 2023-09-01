package internalhttp

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/logger"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/server"
)

const (
	eventsURL = "/v1/calendar/events"
)

type Handler struct {
	app server.Calendar
	log logger.ILogger
}

func NewHandler(log logger.ILogger, app server.Calendar) *Handler {
	return &Handler{
		app: app,
		log: log,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(WithLogger(h.log))

	router.Route(eventsURL, func(r chi.Router) {
		r.Post("/", h.createEvent())
		r.Get("/day", h.getEventsByDay())
		r.Get("/week", h.getEventsByWeek())
		r.Get("/month", h.getEventsByMonth())
		r.Patch("/{id}", h.updateEvent())
		r.Delete("/{id}", h.deleteEvent())
	})
	return router
}
