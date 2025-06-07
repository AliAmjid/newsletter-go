package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load env values
	err := godotenv.Load()
	if err != nil {
		log.Fatal("An error has occured during loading from .env file")
	}
	// env values can be used with os.Getenv("POSTGRES_HOST") command
}
