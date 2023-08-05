package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/storage"
)

type DB interface {
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Close() error
}

type Storage struct {
	db DB
}

func New(db DB) *Storage {
	return &Storage{db: db}
}

const (
	SQLCreateEvent = `
	INSERT INTO events(id, title, description, user_id, start_date, end_date, notification_time)
	VALUES (:id, :title, :description, :user_id, :start_date, :end_date, :notification_time)`

	SQLDeleteEvent = `
	DELETE FROM events
	WHERE id = $1`

	SQLGetByDate = `
	SELECT id, title, description, user_id, start_date, end_date, notification_time FROM events
	WHERE user_id = $1 AND
	(start_date >= $2 AND start_date < $3) OR (end_date >= $4 AND end_date < $5) OR (start_date < $6 AND end_date > $7)`
)

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	_, err := s.db.NamedExecContext(ctx, SQLCreateEvent, event)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	sqlUpdateQuery := buildUpdateQuery(event)
	_, err := s.db.NamedExecContext(ctx, sqlUpdateQuery, event)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	_, err := s.db.ExecContext(ctx, SQLDeleteEvent, eventID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetEventByDay(ctx context.Context, userID int64, day time.Time) ([]storage.Event, error) {
	events, err := s.getEventByDate(ctx, userID, day, day.AddDate(0, 0, 1))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNoEventsFound
		}
		return nil, err
	}

	return events, nil
}

func (s *Storage) GetEventByWeek(ctx context.Context, userID int64, week time.Time) ([]storage.Event, error) {
	return s.getEventByDate(ctx, userID, week, week.AddDate(0, 0, 7))
}

func (s *Storage) GetEventByMonth(ctx context.Context, userID int64, month time.Time) ([]storage.Event, error) {
	return s.getEventByDate(ctx, userID, month, month.AddDate(0, 1, 0))
}

func (s *Storage) getEventByDate(ctx context.Context, userID int64, start, end time.Time) ([]storage.Event, error) {
	var events []storage.Event
	err := s.db.SelectContext(ctx, &events, SQLGetByDate, userID, start, end, start, end, start, end)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func buildUpdateQuery(event storage.Event) string {
	qb := NewUpdateQueryBuilder("events")

	qb.SetIf(event.Title != "", "title = :title")
	qb.SetIf(event.Description != nil, "description = :description")
	qb.SetIf(!event.StartDate.IsZero(), "start_date = :start_date")
	qb.SetIf(!event.EndDate.IsZero(), "end_date = :end_date")
	qb.SetIf(event.NotificationTime != nil, "notification_time = :notification_time")

	qb.Where("id = :id")

	return qb.Build()
}
