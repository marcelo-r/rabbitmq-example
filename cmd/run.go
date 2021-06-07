package cmd

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/marcelo-r/rabbitmq-example/models"
	"github.com/marcelo-r/rabbitmq-example/rabbit"
)

// RabbitURI is the connection url
const RabbitURI = "amqp://guest:guest@localhost:5672/"

// Producer is the entry point for consuming data from a csv file and sending
// each record to a rabbitmq queue
func Produce(delay time.Duration) error {
	filename := "mock_data.csv"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("unable to open file %s, got error: %s", filename, err)
	}
	records := make(chan models.Record)
	go func() {
		log.Println("extraction begun")
		rabbit.ExtractCsvRecords(file, records)
	}()

	mq, err := rabbit.Init(RabbitURI, "example")
	if err != nil {
		log.Fatalf("unable to init rabbitmq: %s", err)
	}
	defer mq.Close()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go rabbit.Publisher(i, mq, records, &wg)
	}
	log.Println("waiting on workers...")
	wg.Wait()
	return nil
}

func Consumer() {
	// open rabbitmq connection

	// open postgres connection

	// consume from queue in bulk

	// insert in bulk
}
