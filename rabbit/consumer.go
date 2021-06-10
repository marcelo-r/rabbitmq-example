package rabbit

import (
	"encoding/json"
	"log"

	"github.com/marcelo-r/rabbitmq-example/models"
	"github.com/streadway/amqp"
)

// Consumer is a worker to receive and serialize messages
func Consumer(name int, input <-chan amqp.Delivery, out chan models.Record) {
	log.Printf("consumer %v: started", name)
	for msg := range input {
		record := models.Record{}
		err := json.Unmarshal(msg.Body, &record)
		if err != nil {
			log.Printf("could not unmarshal: %s", err)
		}
		out <- record
	}
}
