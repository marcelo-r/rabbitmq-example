package producer

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/marcelo-r/rabbitmq-example/models"
)

// ExtractCsvRecords receives a io.Reader and sends back csv records through a channel
func ExtractCsvRecords(r io.Reader) <-chan models.Record {
	out := make(chan models.Record, 100)
	reader := csv.NewReader(r)
	_, _ = reader.Read() // discard header row
	go func() {
		for {
			record, err := reader.Read()
			if record == nil && err == io.EOF {
				break
			} else if err != nil {
				_ = fmt.Errorf("%s", err)
			}
			out <- models.SerializeRecord(record)
		}
		close(out)
	}()
	return out
}
