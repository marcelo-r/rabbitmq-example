package models

import (
	"log"
	"strconv"
	"time"
)

// Record represents the mock data
type Record struct {
	ID          int       `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	Email       string    `json:"email"`
	IP          string    `json:"ip"`
	Mac         string    `json:"mac"`
	CountryCode string    `json:"country_code"`
	UserAgent   string    `json:"user_agent"`
}

// SerializeRecord turns a record in csv format to a Record struct
func SerializeRecord(data []string) Record {
	id, _ := strconv.Atoi(data[0])
	timestamp, err := time.Parse("2006-01-02 15:04:05", data[1])
	if err != nil {
		log.Fatalf("could not parse time from %s: %s", data[1], err)
	}
	return Record{
		ID:          id,
		Timestamp:   timestamp,
		Email:       data[2],
		IP:          data[3],
		Mac:         data[4],
		CountryCode: data[5],
		UserAgent:   data[6],
	}
}
