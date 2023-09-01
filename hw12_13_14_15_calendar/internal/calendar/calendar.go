package calendar

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/models"
)

type Storage interface {
	CreateEvent(context.Context, *models.Event) error
	UpdateEvent(context.Context, *models.Event) error
	DeleteEvent(context.Context, string) error
	GetEventByDay(context.Context, int64, time.Time) ([]models.Event, error)
	GetEventByWeek(context.Context, int64, time.Time) ([]models.Event, error)
	GetEventByMonth(context.Context, int64, time.Time) ([]models.Event, error)
}

type Calendar struct {
	db Storage
}

func New(storage Storage) *Calendar {
	return &Calendar{
		db: storage,
	}
}

func (c *Calendar) CreateEvent(ctx context.Context, event *models.Event) (string, error) {
	event.ID = generateEventID()
	if err := c.db.CreateEvent(ctx, event); err != nil {
		return "", err
	}
	return event.ID, nil
}

func generateEventID() string {
	return uuid.New().String()
}

func (c *Calendar) UpdateEvent(ctx context.Context, event *models.Event) error {
	return c.db.UpdateEvent(ctx, event)
}

func (c *Calendar) DeleteEvent(ctx context.Context, eventID string) error {
	return c.db.DeleteEvent(ctx, eventID)
}

func (c *Calendar) GetEventByDay(ctx context.Context, userID int64, day time.Time) ([]models.Event, error) {
	return c.db.GetEventByDay(ctx, userID, day)
}

func (c *Calendar) GetEventByWeek(ctx context.Context, userID int64, week time.Time) ([]models.Event, error) {
	return c.db.GetEventByWeek(ctx, userID, week)
}

func (c *Calendar) GetEventByMonth(ctx context.Context, userID int64, month time.Time) ([]models.Event, error) {
	return c.db.GetEventByMonth(ctx, userID, month)
}
