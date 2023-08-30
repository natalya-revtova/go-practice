package internalhttp

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/logger"
	"golang.org/x/exp/slog"
)

func WithLogger(log logger.ILogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()

			entry := log.With(
				slog.String("request_id", middleware.GetReqID(r.Context())),
				slog.String("ts", t1.String()),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("protocol_version", r.Proto),
				slog.String("user_agent", r.UserAgent()),
			)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				entry.Info("HTTP/server: request completed",
					slog.Int("status", ww.Status()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
