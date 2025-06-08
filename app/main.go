package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/AliAmjid/newsletter-go/internal/db"
	delivery "github.com/AliAmjid/newsletter-go/internal/delivery/http"
	"github.com/AliAmjid/newsletter-go/internal/repository/postgres"
	postusecase "github.com/AliAmjid/newsletter-go/internal/usecase/post"
)

func main() {
	// Load env values
	if err := godotenv.Load(); err != nil {
		log.Fatal("An error has occured during loading from .env file")
	}

	dbConnectionString := os.Getenv("POSTGRES_CONNECTION_STRING")
	db.Init(dbConnectionString)

	repo := postgres.NewPostRepository(db.DB)
	service := postusecase.NewService(repo)

	r := delivery.NewRouter()
	delivery.NewPostHandler(r, service)

	fmt.Println("Server starting on port 3000")
	server := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("failed to listen to server", err)
	}
}
