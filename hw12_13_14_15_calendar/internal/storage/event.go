package storage

import (
	"context"
	"errors"
	"time"

	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/entity"
)

// ErrEventNotFound ошибка, если событие не найдено.
var ErrEventNotFound = errors.New("event not found")

// ErrInvalidPeriod ошибка, если указан неверный период.
var ErrInvalidPeriod = errors.New("invalid period")

// EventStorage интерфейс хранилища событий.
type EventStorage interface {
	Create(ctx context.Context, event entity.Event) (string, error)
	Update(ctx context.Context, eventID string, event entity.Event) error
	Delete(ctx context.Context, eventID string) error
	List(ctx context.Context, date time.Time, period string) ([]entity.Event, error)
}
