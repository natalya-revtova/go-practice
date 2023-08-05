package storage

import (
	"errors"
	"time"
)

const (
	InMemory = "in-memory"
	SQL      = "sql"
)

var (
	ErrNoEventsFound = errors.New("no events found")
	ErrEventNotExist = errors.New("event does not exist")
)

type Event struct {
	ID               string         `db:"id"`
	Title            string         `db:"title"`
	Description      *string        `db:"description"`
	UserID           int64          `db:"user_id"`
	StartDate        time.Time      `db:"start_date"`
	EndDate          time.Time      `db:"end_date"`
	NotificationTime *time.Duration `db:"notification_time"`
}
