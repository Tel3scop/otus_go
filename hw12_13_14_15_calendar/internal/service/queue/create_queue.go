package queue

import "context"

// CreateQueue создает очередь.
func (s *queueServ) CreateQueue(ctx context.Context) error {
	return s.queueStorage.CreateQueue(ctx)
}
