package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/app"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/config"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/server/http/handlers/hello"
)

type Application interface {
	CreateEvent(context.Context, app.Event) error
	UpdateEvent(context.Context, app.Event) error
	DeleteEvent(context.Context, string) error
	GetEventByDay(context.Context, int64, time.Time) ([]app.Event, error)
	GetEventByWeek(context.Context, int64, time.Time) ([]app.Event, error)
	GetEventByMonth(context.Context, int64, time.Time) ([]app.Event, error)
}

type Server struct {
	srv *http.Server
	app Application
	log app.Logger
}

func NewServer(logger app.Logger, app Application, cfg *config.ServerConfig) *Server {
	serverCfg := http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	server := Server{
		srv: &serverCfg,
		app: app,
		log: logger,
	}

	server.srv.Handler = server.initRoutes()
	return &server
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.srv.ListenAndServe(); err != nil {
		<-ctx.Done()
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) initRoutes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(WithLogger(s.log))

	router.Get("/hello", hello.SayHello())

	return router
}
