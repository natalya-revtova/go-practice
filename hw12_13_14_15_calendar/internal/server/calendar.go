package server

import (
	"context"
	"time"

	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.33.0 --name=Calendar
type Calendar interface {
	CreateEvent(context.Context, *models.Event) error
	UpdateEvent(context.Context, *models.Event) error
	DeleteEvent(context.Context, string) error
	GetEventByDay(context.Context, int64, time.Time) ([]models.Event, error)
	GetEventByWeek(context.Context, int64, time.Time) ([]models.Event, error)
	GetEventByMonth(context.Context, int64, time.Time) ([]models.Event, error)
}
