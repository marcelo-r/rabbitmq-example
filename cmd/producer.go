package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/marcelo-r/rabbitmq-example/producer"
	"github.com/marcelo-r/rabbitmq-example/rabbit"
	"github.com/streadway/amqp"
)

func Producer() {
	filename := "mock_data.csv"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("unable to open file %s, got error: %s", filename, err)
	}
	recordsChan := producer.ExtractCsvRecords(file)

	conn, err := amqp.Dial(rabbit.RabbitURI)
	if err != nil {
		log.Fatalf("could not open connection to rabbitmq: %s", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not open a channel: %s", err)
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare("example", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to declare queue: %s", err)
	}
	for record := range recordsChan {
		err := rabbit.PublishRecord(channel, queue, record)
		if err != nil {
			log.Fatalf("error publishing: %s", err)
		} else {
			go func() {
				fmt.Printf("%d: ok\n", record.ID)
			}()
		}
	}
}
