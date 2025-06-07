package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init(dataSourceName string) {
	var err error
	DB, err = sql.Open("postgres", dataSourceName)

	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to DB")
}
