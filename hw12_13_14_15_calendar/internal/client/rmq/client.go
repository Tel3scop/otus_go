package rmq

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

// Client представляет клиент RabbitMQ.
type Client struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

// NewClient создает новый клиент RabbitMQ.
func NewClient(url string) (*Client, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &Client{conn: conn, ch: ch}, nil
}

// CreateQueue создает очередь.
func (c *Client) CreateQueue(queueName string) error {
	_, err := c.ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}
	return nil
}

// Publish помещает сообщение в очередь.
func (c *Client) Publish(queueName string, body []byte) error {
	err := c.ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}
	return nil
}

// Consume читает сообщение из очереди.
func (c *Client) Consume(queueName string) (<-chan amqp091.Delivery, error) {
	msgs, err := c.ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %w", err)
	}
	return msgs, nil
}

// DeleteQueue удаляет очередь.
func (c *Client) DeleteQueue(queueName string) error {
	_, err := c.ch.QueueDelete(
		queueName, // name
		false,     // ifUnused
		false,     // ifEmpty
		false,     // no-wait
	)
	if err != nil {
		return fmt.Errorf("failed to delete a queue: %w", err)
	}
	return nil
}

// Close закрывает соединение и канал.
func (c *Client) Close() {
	c.ch.Close()
	c.conn.Close()
}
