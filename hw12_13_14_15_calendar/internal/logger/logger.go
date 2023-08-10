package logger

import (
	"os"

	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/app"
	"golang.org/x/exp/slog"
)

type Logger struct {
	log *slog.Logger
}

func New(level slog.Level) Logger {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	return Logger{log: log}
}

func (l Logger) With(args ...any) app.Logger {
	l.log = l.log.With(args...)
	return l
}

func (l Logger) Info(msg string, args ...any) {
	l.log.Info(msg, args...)
}

func (l Logger) Debug(msg string, args ...any) {
	l.log.Debug(msg, args...)
}

func (l Logger) Warn(msg string, args ...any) {
	l.log.Warn(msg, args...)
}

func (l Logger) Error(msg string, args ...any) {
	l.log.Error(msg, args...)
}
