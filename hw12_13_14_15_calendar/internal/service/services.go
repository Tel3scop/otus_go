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
	DeleteByDate(ctx context.Context, date time.Time) error
	List(ctx context.Context, date time.Time, period entity.PeriodType) ([]entity.Event, error)
}

// QueueService интерфейс сервиса для работы с очередями.
type QueueService interface {
	CreateQueue(ctx context.Context) error
	Enqueue(ctx context.Context, message []byte) error
	Dequeue(ctx context.Context) ([]byte, error)
	DeleteQueue(ctx context.Context) error
}

// NotificationService интерфейс сервиса для работы с уведомлениями.
type NotificationService interface {
	Enqueue(ctx context.Context, notification entity.Notification) error
	Dequeue(ctx context.Context) (*entity.Notification, error)
	NotifyOnEvent(ctx context.Context) error
	Send(ctx context.Context) error
}
