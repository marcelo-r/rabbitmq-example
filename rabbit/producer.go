package rabbit

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/marcelo-r/rabbitmq-example/models"
)

// ExtractCsvRecords receives a io.Reader and sends back csv records through a channel
func ExtractCsvRecords(r io.Reader, out chan models.Record) int {
	reader := csv.NewReader(r)
	_, _ = reader.Read() // discard header row
	i := 1
	go func() {
		for ; ; i++ {
			record, err := reader.Read()
			if record == nil && err == io.EOF {
				break
			} else if err != nil {
				_ = fmt.Errorf("%s", err)
			}
			out <- models.SerializeRecord(record)
		}
		close(out)
		log.Printf("closing records: %d", len(out))
	}()
	return i
}

func Publisher(name int, mq *RabbitMQ, input chan models.Record, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("publisher %d: started", name)

	for r := range input {
		err := mq.PublishRecord(r)
		if err != nil {
			log.Printf("publisher %d can't PublishRecord: %s", name, err)
		}
	}
	log.Printf("publisher %d: done", name)
}
