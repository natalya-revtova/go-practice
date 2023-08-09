package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/storage"
)

type Logger interface {
	With(args ...any) Logger
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

type Storage interface {
	CreateEvent(context.Context, storage.Event) error
	UpdateEvent(context.Context, storage.Event) error
	DeleteEvent(context.Context, string) error
	GetEventByDay(context.Context, int64, time.Time) ([]storage.Event, error)
	GetEventByWeek(context.Context, int64, time.Time) ([]storage.Event, error)
	GetEventByMonth(context.Context, int64, time.Time) ([]storage.Event, error)
}

type App struct {
	storage Storage
	log     Logger
}

func New(logger Logger, storage Storage) *App {
	return &App{
		storage: storage,
		log:     logger,
	}
}

func (a *App) CreateEvent(ctx context.Context, event Event) error {
	event.ID = generateEventID()
	if err := a.storage.CreateEvent(ctx, event.ToStorageModel()); err != nil {
		a.log.Error("Can not add event to storage", "event_id", event.ID, "error", err)
		return err
	}
	return nil
}

func generateEventID() string {
	return uuid.New().String()
}

func (a *App) UpdateEvent(ctx context.Context, event Event) error {
	if err := a.storage.UpdateEvent(ctx, event.ToStorageModel()); err != nil {
		a.log.Error("Can not update event", "event_id", event.ID, "error", err)
		return err
	}
	return nil
}

func (a *App) DeleteEvent(ctx context.Context, eventID string) error {
	if err := a.storage.DeleteEvent(ctx, eventID); err != nil {
		a.log.Error("Can not delete event", "event_id", eventID, "error", err)
		return err
	}
	return nil
}

func (a *App) GetEventByDay(ctx context.Context, userID int64, day time.Time) ([]Event, error) {
	events, err := a.storage.GetEventByDay(ctx, userID, day)
	if err != nil {
		a.log.Error("Can not get events for the selected day", "user_id", userID, "day", day, "error", err)
		return nil, err
	}
	result := make(Events, len(events))
	return result.FromStorageModel(events), nil
}

func (a *App) GetEventByWeek(ctx context.Context, userID int64, week time.Time) ([]Event, error) {
	events, err := a.storage.GetEventByWeek(ctx, userID, week)
	if err != nil {
		a.log.Error("Can not get events for the selected week", "user_id", userID, "week", week, "error", err)
		return nil, err
	}
	result := make(Events, len(events))
	return result.FromStorageModel(events), nil
}

func (a *App) GetEventByMonth(ctx context.Context, userID int64, month time.Time) ([]Event, error) {
	events, err := a.storage.GetEventByMonth(ctx, userID, month)
	if err != nil {
		a.log.Error("Can not get events for the selected month", "user_id", userID, "month", month, "error", err)
		return nil, err
	}
	result := make(Events, len(events))
	return result.FromStorageModel(events), nil
}
