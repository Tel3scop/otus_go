package rmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tel3scop/helpers/logger"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/client/rmq"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/entity"
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

// Enqueue помещает уведомление в очередь.
func (r *queueRepo) Enqueue(_ context.Context, notification entity.Notification) error {
	body, err := json.Marshal(notification)
	if err != nil {
		logger.Error("failed to encode notification", zap.Error(err), zap.Any("notification", notification))
		return fmt.Errorf("failed to marshal notification: %v", err)
	}
	return r.client.Publish(queueName, body)
}

// Dequeue читает уведомление из очереди.
func (r *queueRepo) Dequeue(ctx context.Context) (*entity.Notification, error) {
	messages, err := r.client.Consume(queueName)
	if err != nil {
		logger.Error("failed to consume messages", zap.Error(err))
		return nil, fmt.Errorf("failed to consume messages: %v", err)
	}

	select {
	case msg := <-messages:
		var notification entity.Notification
		err := json.Unmarshal(msg.Body, &notification)
		if err != nil {
			logger.Error("failed to unmarshal notification", zap.Error(err), zap.Any("notification", notification))
			return nil, fmt.Errorf("failed to unmarshal notification: %v", err)
		}
		return &notification, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// DeleteQueue удаляет очередь.
func (r *queueRepo) DeleteQueue(_ context.Context) error {
	return r.client.DeleteQueue(queueName)
}
