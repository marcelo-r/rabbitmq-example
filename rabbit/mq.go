package rabbit

import (
	"encoding/json"
	"fmt"

	"github.com/marcelo-r/rabbitmq-example/models"
	"github.com/streadway/amqp"
)

// RabbitURI contains the connection url
const ContentType = "text/json"

// Rabbit contains rabbitmq params
type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   amqp.Queue
}

// Init setups RabbitMQ
func Init(url string, queueName string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("could not open connection to rabbitmq: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("could not open a channel: %w", err)
	}

	queue, err := channel.QueueDeclare("example", false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}
	return &RabbitMQ{
		Conn:    conn,
		Channel: channel,
		Queue:   queue,
	}, nil
}

// Close RabbitMQ's amqp.Connection and amqp.Channel
func (r *RabbitMQ) Close() {
	r.Conn.Close()
	r.Channel.Close()
}

// PublishRecord sends a record to queue with name equal to queue on channel
func (r *RabbitMQ) PublishRecord(record models.Record) error {
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("could not marshal record to json: %w", err)
	}
	err = r.Channel.Publish(
		"",
		r.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: ContentType,
			Body:        recordJSON,
		},
	)
	if err != nil {
		r.Channel.Close()
		r.Channel = nil // avoid calling erroed channel, amqp.Channel requires recreate
		return fmt.Errorf("could not publish record to %s: %w", r.Queue.Name, err)
	}
	return nil
}
