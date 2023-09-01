package internalhttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/config"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/logger"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/server"
)

type Server struct {
	srv *http.Server
}

func NewServer(log logger.ILogger, app server.Calendar, cfg *config.ServerHTTPConfig) *Server {
	serverCfg := http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      NewHandler(log, app).InitRoutes(),
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		srv: &serverCfg,
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
