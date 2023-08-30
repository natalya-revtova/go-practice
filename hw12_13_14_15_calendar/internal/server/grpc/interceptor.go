package grpc

import (
	"context"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/logger"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

func UnaryLoggerInterceptor(log logger.ILogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		t1 := time.Now()
		reqID := newRequestID()

		entry := log.With(
			slog.String("request_id", reqID),
			slog.String("method", info.FullMethod),
		)

		defer func() {
			entry.Info("gRPC/server: request completed",
				slog.String("duration", time.Since(t1).String()),
			)
		}()

		ctx = context.WithValue(ctx, middleware.RequestIDKey, reqID)
		return handler(ctx, req)
	}
}

func newRequestID() string {
	return uuid.New().String()
}
