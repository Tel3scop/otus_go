package rmq

import (
	"context"
	"fmt"

	"github.com/Tel3scop/helpers/logger"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/client/rmq"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/storage"
	"go.uber.org/zap"
)

const (
	queueName = "notification_queue"
)

type queueRepo struct {
	client *rmq.Client
}

// NewRepository создает новое хранилище для работы с очередями.
func NewRepository(client *rmq.Client) storage.QueueStorage {
	return &queueRepo{client: client}
}

// CreateQueue создает очередь.
func (r *queueRepo) CreateQueue(_ context.Context) error {
	return r.client.CreateQueue(queueName)
}

// Enqueue помещает сообщение в очередь.
func (r *queueRepo) Enqueue(_ context.Context, message []byte) error {
	return r.client.Publish(queueName, message)
}

// Dequeue читает сообщение из очереди.
func (r *queueRepo) Dequeue(ctx context.Context) ([]byte, error) {
	messages, err := r.client.Consume(queueName)
	if err != nil {
		logger.Error("failed to consume messages", zap.Error(err))
		return nil, fmt.Errorf("failed to consume messages: %w", err)
	}

	select {
	case msg := <-messages:
		return msg.Body, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// DeleteQueue удаляет очередь.
func (r *queueRepo) DeleteQueue(_ context.Context) error {
	return r.client.DeleteQueue(queueName)
}
