package grpc

import (
	"fmt"
	"net"

	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/config"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/logger"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/server"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/pkg/api/calendarpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	srv *grpc.Server
	calendarpb.UnimplementedCalendarServer
	app server.Calendar
	log logger.ILogger
}

func NewServer(logger logger.ILogger, app server.Calendar, cfg *config.ServerGRPCConfig) *Server {
	var serverOptions []grpc.ServerOption
	if cfg != nil {
		serverOptions = []grpc.ServerOption{
			grpc.Creds(insecure.NewCredentials()),
			grpc.UnaryInterceptor(UnaryLoggerInterceptor(logger)),
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle: cfg.MaxConnectionIdle,
				MaxConnectionAge:  cfg.MaxConnectionAge,
				Time:              cfg.Time,
				Timeout:           cfg.Timeout,
			}),
		}
	}

	srv := grpc.NewServer(serverOptions...)

	return &Server{
		srv: srv,
		app: app,
		log: logger,
	}
}

func (s *Server) Start(cfg *config.ServerGRPCConfig) error {
	lsn, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		return err
	}

	calendarpb.RegisterCalendarServer(s.srv, s)
	return s.srv.Serve(lsn)
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
