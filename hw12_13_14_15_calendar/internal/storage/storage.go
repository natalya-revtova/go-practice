package storage

import (
	"errors"
)

const (
	InMemory = "in-memory"
	SQL      = "sql"
)

var (
	ErrNoEventsFound = errors.New("no events found")
	ErrEventNotExist = errors.New("event does not exist")
)
