package queue

import (
	"context"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/entity"
)

// DequeueNotification читает уведомление из очереди.
func (s *queueServ) DequeueNotification(ctx context.Context) (*entity.Notification, error) {
	return s.queueStorage.Dequeue(ctx)
}
