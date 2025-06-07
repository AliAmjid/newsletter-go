package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AliAmjid/newsletter-go/internal/db"
	"github.com/AliAmjid/newsletter-go/internal/router"
	"github.com/joho/godotenv"
)

func main() {
	// Load env values
	err := godotenv.Load()
	if err != nil {
		log.Fatal("An error has occured during loading from .env file")
	}

	dbConnectionString := os.Getenv("POSTGRES_CONNECTION_STRING")
	db.Init(dbConnectionString)

	startServer()
}

func startServer() {
	fmt.Println("Server starting on port 3000")

	router := router.NewRouter()

	server := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("failed to listen to server", err)
	}
}
