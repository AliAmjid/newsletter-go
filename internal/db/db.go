package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init(dataSourceName string) {
	var err error
	for {
		DB, err = sql.Open("postgres", dataSourceName)
		if err != nil {
			log.Printf("Failed to open database: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		if err = DB.Ping(); err != nil {
			log.Printf("Failed to connect to database: %v. Retrying in 5 seconds...", err)
			DB.Close()
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	log.Println("Connected to DB")
}
