package service

import (
	"context"
	"time"

	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/entity"
)

// EventService интерфейс сервиса событий.
type EventService interface {
	Create(ctx context.Context, event entity.Event) (string, error)
	Update(ctx context.Context, event entity.Event) error
	Delete(ctx context.Context, eventID string) error
	List(ctx context.Context, date time.Time, period entity.PeriodType) ([]entity.Event, error)
}
