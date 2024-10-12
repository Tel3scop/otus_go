package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/client/db"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/entity"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type NoOpTxManager struct{}

func (n *NoOpTxManager) ReadCommitted(ctx context.Context, f db.Handler) error {
	return f(ctx)
}

// InMemEventStorage реализация хранилища событий в памяти.
type InMemEventStorage struct {
	events map[string]entity.Event
	mu     sync.RWMutex
}

// NewInMemoryEventStorage создание нового хранилища событий в памяти.
func NewInMemoryEventStorage() *InMemEventStorage {
	return &InMemEventStorage{
		events: make(map[string]entity.Event),
	}
}

// Create добавление события в хранилище.
func (s *InMemEventStorage) Create(_ context.Context, event entity.Event) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	event.ID = uuid.New().String()
	s.events[event.ID] = event
	return event.ID, nil
}

// Update обновление события в хранилище.
func (s *InMemEventStorage) Update(_ context.Context, eventID string, event entity.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.events[eventID]; !exists {
		return storage.ErrEventNotFound
	}
	event.ID = eventID
	s.events[eventID] = event
	return nil
}

// Delete удаление события из хранилища.
func (s *InMemEventStorage) Delete(_ context.Context, eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.events[eventID]; !exists {
		return storage.ErrEventNotFound
	}
	delete(s.events, eventID)
	return nil
}

// DeleteByDate удаление событий из хранилища по дате.
func (s *InMemEventStorage) DeleteByDate(_ context.Context, date time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, event := range s.events {
		if event.DateTime.Before(date) {
			delete(s.events, event.ID)
		}
	}

	return nil
}

// List список событий на определенный период.
func (s *InMemEventStorage) List(_ context.Context, date time.Time, period entity.PeriodType) ([]entity.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var events []entity.Event
	for _, event := range s.events {
		switch period {
		case entity.PeriodDay:
			if event.DateTime.Year() == date.Year() &&
				event.DateTime.Month() == date.Month() &&
				event.DateTime.Day() == date.Day() {
				events = append(events, event)
			}
		case entity.PeriodWeek:
			startOfWeek := date.AddDate(0, 0, -int(date.Weekday()))
			endOfWeek := startOfWeek.AddDate(0, 0, 7)
			if event.DateTime.After(startOfWeek) && event.DateTime.Before(endOfWeek) {
				events = append(events, event)
			}
		case entity.PeriodMonth:
			startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
			endOfMonth := startOfMonth.AddDate(0, 1, 0)
			if event.DateTime.After(startOfMonth) && event.DateTime.Before(endOfMonth) {
				events = append(events, event)
			}
		default:
			return nil, storage.ErrInvalidPeriod
		}
	}
	return events, nil
}
