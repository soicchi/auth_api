package main

import (
	"log"

	"github.com/soicchi/auth_api/models"
)

func main() {
	// Create database connection
	_, err := models.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	log.Println("Successfully connected to database")
}
