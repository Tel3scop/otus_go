package queue

import (
	"context"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/entity"
)

// EnqueueNotification помещает уведомление в очередь.
func (s *queueServ) EnqueueNotification(ctx context.Context, notification entity.Notification) error {
	return s.queueStorage.Enqueue(ctx, notification)
}
