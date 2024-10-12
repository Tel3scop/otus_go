package queue

import "context"

// DeleteQueue удаляет очередь.
func (s *queueServ) DeleteQueue(ctx context.Context) error {
	return s.queueStorage.DeleteQueue(ctx)
}
