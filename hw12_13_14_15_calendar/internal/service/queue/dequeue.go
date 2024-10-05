package queue

import (
	"context"
)

// Dequeue читает сообщение из очереди.
func (s *queueServ) Dequeue(ctx context.Context) ([]byte, error) {
	return s.queueStorage.Dequeue(ctx)
}
