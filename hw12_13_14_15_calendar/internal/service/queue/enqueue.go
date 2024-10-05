package queue

import (
	"context"
)

// Enqueue помещает сообщение в очередь.
func (s *queueServ) Enqueue(ctx context.Context, message []byte) error {
	return s.queueStorage.Enqueue(ctx, message)
}
