package memorystorage

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/models"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/storage"
)

type (
	id = string

	events map[id]*models.Event
	dates  map[time.Time]map[id]struct{}
)

type Storage struct {
	events events
	days   dates
	weeks  dates
	months dates
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{
		events: make(events),
		days:   make(dates),
		weeks:  make(dates),
		months: make(dates),
	}
}

func (s *Storage) CreateEvent(ctx context.Context, event *models.Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[event.ID] = event
	s.saveDates(event.ID, event.Day, event.Week, event.Month)

	return nil
}

func (s *Storage) saveDates(eventID string, day, week, month time.Time) {
	if _, ok := s.days[day]; !ok {
		s.days[day] = make(map[id]struct{})
	}
	s.days[day][eventID] = struct{}{}

	if _, ok := s.weeks[week]; !ok {
		s.weeks[week] = make(map[id]struct{})
	}
	s.weeks[week][eventID] = struct{}{}

	if _, ok := s.months[month]; !ok {
		s.months[month] = make(map[id]struct{})
	}
	s.months[month][eventID] = struct{}{}
}

func (s *Storage) UpdateEvent(ctx context.Context, event *models.Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	updated, ok := s.events[event.ID]
	if !ok {
		return storage.ErrEventNotExist
	}

	updateEventFields(updated, event)
	s.events[event.ID] = updated

	return nil
}

func updateEventFields(updated *models.Event, event *models.Event) {
	if event.Description != nil {
		updated.Description = event.Description
	}
	if len(event.Title) != 0 {
		updated.Title = event.Title
	}
	if !event.StartDate.IsZero() {
		updated.StartDate = event.StartDate
	}
	if !event.Day.IsZero() {
		updated.Day = event.Day
	}
	if !event.Week.IsZero() {
		updated.Week = event.Week
	}
	if !event.Month.IsZero() {
		updated.Month = event.Month
	}
	if !event.EndDate.IsZero() {
		updated.EndDate = event.EndDate
	}
	if event.NotificationTime != nil {
		updated.NotificationTime = event.NotificationTime
	}
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	deleted, ok := s.events[eventID]
	if !ok {
		return storage.ErrEventNotExist
	}

	s.deleteDates(eventID, deleted.Day, deleted.Week, deleted.Month)
	delete(s.events, eventID)

	return nil
}

func (s *Storage) deleteDates(eventID string, day, week, month time.Time) {
	delete(s.days[day], eventID)
	if len(s.days[day]) == 0 {
		delete(s.days, day)
	}

	delete(s.weeks[week], eventID)
	if len(s.weeks[week]) == 0 {
		delete(s.weeks, week)
	}

	delete(s.months[month], eventID)
	if len(s.months[month]) == 0 {
		delete(s.months, month)
	}
}

func (s *Storage) GetEventByDay(ctx context.Context, userID int64, day time.Time) ([]models.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getSortedEventsByIDs(userID, s.days[day]), nil
}

func (s *Storage) GetEventByWeek(ctx context.Context, userID int64, week time.Time) ([]models.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getSortedEventsByIDs(userID, s.weeks[week]), nil
}

func (s *Storage) GetEventByMonth(ctx context.Context, userID int64, month time.Time) ([]models.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getSortedEventsByIDs(userID, s.months[month]), nil
}

func (s *Storage) getSortedEventsByIDs(userID int64, ids map[id]struct{}) []models.Event {
	if len(ids) == 0 {
		return nil
	}

	events := make([]models.Event, 0, len(ids))
	for id := range ids {
		if userID == s.events[id].UserID {
			events = append(events, *s.events[id])
		}
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].StartDate.Before(events[j].StartDate)
	})
	return events
}
