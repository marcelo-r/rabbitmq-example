package rabbit

import (
	"encoding/json"
	"fmt"

	"github.com/marcelo-r/rabbitmq-example/models"
	"github.com/streadway/amqp"
)

// RabbitURI contains the connection url
const RabbitURI = "amqp://guest:guest@localhost:5672/"
const ContentType = "text/json"

// PublishRecord sends a record to queue with name equal to queue on channel
func PublishRecord(channel *amqp.Channel, queue amqp.Queue, record models.Record) error {
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("could not marshal record to json: %w", err)
	}
	err = channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: ContentType,
			Body:        recordJSON,
		},
	)
	if err != nil {
		return fmt.Errorf("could not publish record to %s: %w", queue.Name, err)
	}
	return nil
}
