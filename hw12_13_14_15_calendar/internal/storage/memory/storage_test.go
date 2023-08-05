package memorystorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestCreateEvent(t *testing.T) {
	memoryStorage := New()

	newEvents := generateEvents(time.Date(2010, 1, 1, 1, 13, 0, 0, time.Local),
		time.Date(2010, 1, 1, 1, 13, 0, 0, time.Local), 1)
	wantEvents := getEvents(newEvents)
	wantDates := getDates(newEvents)

	err := memoryStorage.CreateEvent(context.Background(), newEvents[0])
	require.NoError(t, err)
	require.Equal(t, wantEvents, memoryStorage.events)
	require.Equal(t, wantDates, memoryStorage.dates)
}

func TestDeleteEvent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		memoryStorage := New()

		newEvents := generateEvents(time.Date(2010, 1, 1, 1, 13, 0, 0, time.Local),
			time.Date(2010, 1, 1, 1, 13, 0, 0, time.Local), 1)
		wantEvents := getEvents(nil)
		wantDates := getDates(nil)

		err := memoryStorage.CreateEvent(context.Background(), newEvents[0])
		require.NoError(t, err)

		err = memoryStorage.DeleteEvent(context.Background(), newEvents[0].ID)
		require.NoError(t, err)
		require.Equal(t, wantEvents, memoryStorage.events)
		require.Equal(t, wantDates, memoryStorage.dates)
	})

	t.Run("event does not exist", func(t *testing.T) {
		memoryStorage := New()

		err := memoryStorage.DeleteEvent(context.Background(), "id")
		require.ErrorIs(t, err, storage.ErrEventNotExist)
	})
}

func TestUpdateEvent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		memoryStorage := New()

		newEvents := generateEvents(time.Date(2010, 1, 1, 1, 13, 0, 0, time.Local),
			time.Date(2010, 1, 1, 1, 13, 0, 0, time.Local), 1)

		err := memoryStorage.CreateEvent(context.Background(), newEvents[0])
		require.NoError(t, err)
		require.Equal(t, newEvents[0].Title, memoryStorage.events[newEvents[0].ID].Title)

		updatedEvent := storage.Event{
			ID:    newEvents[0].ID,
			Title: "new title",
		}

		err = memoryStorage.UpdateEvent(context.Background(), updatedEvent)
		require.NoError(t, err)
		require.Equal(t, updatedEvent.Title, memoryStorage.events[newEvents[0].ID].Title)
	})

	t.Run("event does not exist", func(t *testing.T) {
		memoryStorage := New()

		err := memoryStorage.UpdateEvent(context.Background(), storage.Event{ID: "id"})
		require.ErrorIs(t, err, storage.ErrEventNotExist)
	})
}

func TestGetEventByDay(t *testing.T) {
	memoryStorage := New()

	newEvents := generateEvents(time.Date(2010, 1, 1, 1, 13, 0, 0, time.Local),
		time.Date(2010, 1, 1, 1, 13, 0, 0, time.Local), 2)

	err := memoryStorage.CreateEvent(context.Background(), newEvents[0])
	require.NoError(t, err)

	err = memoryStorage.CreateEvent(context.Background(), newEvents[1])
	require.NoError(t, err)

	got, err := memoryStorage.GetEventByDay(context.Background(), newEvents[0].UserID, getDate(newEvents[0].StartDate))
	require.NoError(t, err)
	require.Equal(t, []storage.Event{newEvents[0]}, got)
}

func TestGetEventByWeek(t *testing.T) {
	memoryStorage := New()

	newEvents := generateEvents(time.Date(2010, 1, 1, 1, 13, 0, 0, time.Local),
		time.Date(2010, 1, 1, 1, 13, 0, 0, time.Local), 2)

	err := memoryStorage.CreateEvent(context.Background(), newEvents[0])
	require.NoError(t, err)

	err = memoryStorage.CreateEvent(context.Background(), newEvents[1])
	require.NoError(t, err)

	got, err := memoryStorage.GetEventByWeek(context.Background(), newEvents[0].UserID, getDate(newEvents[0].StartDate))
	require.NoError(t, err)

	require.Equal(t, newEvents, got)
}

func TestGetEventByMonth(t *testing.T) {
	memoryStorage := New()

	newEvents := generateEvents(time.Date(2010, 1, 1, 1, 13, 0, 0, time.Local),
		time.Date(2010, 1, 1, 1, 13, 0, 0, time.Local), 2)

	err := memoryStorage.CreateEvent(context.Background(), newEvents[0])
	require.NoError(t, err)

	err = memoryStorage.CreateEvent(context.Background(), newEvents[1])
	require.NoError(t, err)

	got, err := memoryStorage.GetEventByMonth(context.Background(), newEvents[0].UserID, getDate(newEvents[0].StartDate))
	require.NoError(t, err)

	require.Equal(t, newEvents, got)
}

func TestStorage(t *testing.T) {
	t.Run("get event that take more than one day", func(t *testing.T) {
		memoryStorage := New()
		newEvents := generateEvents(time.Date(2023, 7, 30, 12, 0, 0, 0, time.Local),
			time.Date(2023, 7, 31, 14, 0, 0, 0, time.Local), 1)

		err := memoryStorage.CreateEvent(context.Background(), newEvents[0])
		require.NoError(t, err)

		fmt.Println(memoryStorage.events)
		fmt.Println(memoryStorage.dates)

		got, err := memoryStorage.GetEventByDay(context.Background(), newEvents[0].UserID,
			time.Date(2023, 7, 30, 0, 0, 0, 0, time.Local))

		require.NoError(t, err)
		require.Equal(t, newEvents, got)

		got, err = memoryStorage.GetEventByDay(context.Background(), newEvents[0].UserID,
			time.Date(2023, 7, 31, 0, 0, 0, 0, time.Local))
		require.NoError(t, err)
		require.Equal(t, newEvents, got)

		got, err = memoryStorage.GetEventByWeek(context.Background(), newEvents[0].UserID,
			time.Date(2023, 7, 28, 0, 0, 0, 0, time.Local))
		require.NoError(t, err)
		require.Equal(t, newEvents, got)

		got, err = memoryStorage.GetEventByMonth(context.Background(), newEvents[0].UserID,
			time.Date(2023, 7, 28, 0, 0, 0, 0, time.Local))
		require.NoError(t, err)
		require.Equal(t, newEvents, got)
	})
}

func getDates(events []storage.Event) Dates {
	dates := make(Dates)
	for i := range events {
		startDate := getDate(events[i].StartDate)
		dates[startDate] = make(map[ID]struct{})
		dates[startDate][events[i].ID] = struct{}{}
	}

	return dates
}

func getEvents(events []storage.Event) Events {
	result := make(Events)
	for i := range events {
		result[events[i].ID] = events[i]
	}

	return result
}

func generateEvents(start, end time.Time, count int) []storage.Event {
	events := make([]storage.Event, 0, count)

	for i := 0; i < count; i++ {
		events = append(events, storage.Event{
			ID:               uuid.New().String(),
			Title:            "some title",
			Description:      nil,
			UserID:           1,
			StartDate:        start,
			EndDate:          end,
			NotificationTime: nil,
		})

		start = start.AddDate(0, 0, 1)
		end = end.AddDate(0, 0, 1)
	}

	return events
}
