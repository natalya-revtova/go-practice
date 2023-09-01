package models

import (
	"time"
)

type Event struct {
	ID               string         `db:"id"`
	Title            string         `db:"title"`
	Description      *string        `db:"description"`
	UserID           int64          `db:"user_id"`
	StartDate        time.Time      `db:"start_date"`
	EndDate          time.Time      `db:"end_date"`
	NotificationTime *time.Duration `db:"notification_time"`

	Day   time.Time `db:"day"`
	Week  time.Time `db:"week"`
	Month time.Time `db:"month"`
}
