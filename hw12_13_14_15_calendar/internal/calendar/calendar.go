package calendar

import (
	"context"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/logger"
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
	db  Storage
	log logger.ILogger
}

func New(logger logger.ILogger, storage Storage) *Calendar {
	return &Calendar{
		db:  storage,
		log: logger,
	}
}

func (c *Calendar) CreateEvent(ctx context.Context, event *models.Event) error {
	event.ID = generateEventID()
	if err := c.db.CreateEvent(ctx, event.UpdateForStorage()); err != nil {
		c.log.Error("Can not add event to storage",
			"request_id", middleware.GetReqID(ctx),
			"event_id", event.ID,
			"error", err)
		return err
	}
	return nil
}

func generateEventID() string {
	return uuid.New().String()
}

func (c *Calendar) UpdateEvent(ctx context.Context, event *models.Event) error {
	if err := c.db.UpdateEvent(ctx, event.UpdateForStorage()); err != nil {
		c.log.Error("Can not update event",
			"request_id", middleware.GetReqID(ctx),
			"event_id", event.ID,
			"error", err)
		return err
	}
	return nil
}

func (c *Calendar) DeleteEvent(ctx context.Context, eventID string) error {
	if err := c.db.DeleteEvent(ctx, eventID); err != nil {
		c.log.Error("Can not delete event",
			"request_id", middleware.GetReqID(ctx),
			"event_id", eventID,
			"error", err)
		return err
	}
	return nil
}

func (c *Calendar) GetEventByDay(ctx context.Context, userID int64, day time.Time) ([]models.Event, error) {
	events, err := c.db.GetEventByDay(ctx, userID, day)
	if err != nil {
		c.log.Error("Can not get events for the selected day",
			"request_id", middleware.GetReqID(ctx),
			"user_id", userID,
			"day", day,
			"error", err)
		return nil, err
	}
	return events, nil
}

func (c *Calendar) GetEventByWeek(ctx context.Context, userID int64, week time.Time) ([]models.Event, error) {
	events, err := c.db.GetEventByWeek(ctx, userID, week)
	if err != nil {
		c.log.Error("Can not get events for the selected week",
			"request_id", middleware.GetReqID(ctx),
			"user_id", userID,
			"week", week,
			"error", err)
		return nil, err
	}
	return events, nil
}

func (c *Calendar) GetEventByMonth(ctx context.Context, userID int64, month time.Time) ([]models.Event, error) {
	events, err := c.db.GetEventByMonth(ctx, userID, month)
	if err != nil {
		c.log.Error("Can not get events for the selected month",
			"request_id", middleware.GetReqID(ctx),
			"user_id", userID,
			"month", month,
			"error", err)
		return nil, err
	}
	return events, nil
}
