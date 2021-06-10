package cmd

import (
	"database/sql"
	"log"
	"os"
	"sync"
	"time"

	"github.com/marcelo-r/rabbitmq-example/models"
	"github.com/marcelo-r/rabbitmq-example/rabbit"
)

// RabbitURI is the connection url
const (
	RabbitURI = "amqp://guest:guest@localhost:5672/"
	PgConnStr = "postgres://postgres:secretkey@localhost/rabbit?sslmode=disable"
	WORKERS   = 20
	BatchSize = 1000
)

// Produce is the entry point for consuming data from a csv file and sending
// each record to a rabbitmq queue
func Produce(delay time.Duration, filename string) error {
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
	for i := 0; i < WORKERS; i++ {
		wg.Add(1)
		go rabbit.Publisher(i, mq, records, &wg)
	}
	log.Println("waiting on workers...")
	wg.Wait()
	return nil
}

// Consume waits to process messages
func Consume() {
	// open rabbitmq connection
	mq, err := rabbit.Init(RabbitURI, "example")
	if err != nil {
		log.Fatalf("can't init rabbitmq: %s", err)
	}
	// open postgres connection
	db, err := sql.Open("postgres", PgConnStr)
	if err != nil {
		log.Fatalf("can't connect to postgres: %s", err)
	}
	// consume from queue in bulk
	messages, err := mq.Channel.Consume(
		mq.Queue.Name,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalln("unable to consume from channel:", err)
	}
	var wg sync.WaitGroup
	records := make(chan models.Record)
	for i := 0; i < WORKERS; i++ {
		wg.Add(1)
		go rabbit.Consumer(i, messages, records)
	}

	var batch []models.Record
	count := 0
	for r := range records {
		batch = append(batch, r)
		if len(batch) == BatchSize {
			count++
			go models.BulkInsertTo(count, db, batch)
			batch = nil
		}
	}
}
