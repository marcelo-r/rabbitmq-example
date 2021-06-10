package models

import (
	"database/sql"
	"log"
	"time"

	"github.com/lib/pq"
)

// Record represents the mock data
type Record struct {
	ID          string    `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	Email       string    `json:"email"`
	IP          string    `json:"ip"`
	Mac         string    `json:"mac"`
	CountryCode string    `json:"country_code"`
	UserAgent   string    `json:"user_agent"`
}

// SerializeRecord turns a record in csv format to a Record struct
func SerializeRecord(data []string) Record {
	timestamp, err := time.Parse("2006-01-02 15:04:05", data[1])
	if err != nil {
		log.Fatalf("could not parse time from %s: %s", data[1], err)
	}
	return Record{
		ID:          data[0],
		Timestamp:   timestamp,
		Email:       data[2],
		IP:          data[3],
		Mac:         data[4],
		CountryCode: data[5],
		UserAgent:   data[6],
	}
}

// BulkInsertTo inserts records in bulk to a sql.DB
func BulkInsertTo(name int, db *sql.DB, records []Record) {
	log.Printf("new batch: %v", name)
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln("could not init transaction:", err)
	}
	stmt, err := tx.Prepare(pq.CopyIn("record", "seq_id", "datetime", "email", "ipv4", "mac", "country_code", "user_agent"))
	if err != nil {
		log.Fatalln("could not prepare statment")
	}
	for _, r := range records {
		_, err := stmt.Exec(r.ID, r.Timestamp, r.Email, r.IP, r.Mac, r.CountryCode, r.UserAgent)
		if err != nil {
			log.Fatalf("could not")
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatalln("unable to flush data:", err)
	}
	err = stmt.Close()
	if err != nil {
		log.Fatalln("unable to close prepared statment:", err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalln("unable to commit:", err)
	}
	log.Printf("end batch: %v", name)
}
