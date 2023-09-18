package main

import (
	"log"

	"github.com/soicchi/auth_api/models"
	"github.com/soicchi/auth_api/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	// Create database connection
	_, err := models.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	log.Println("Successfully connected to database")

	e := echo.New()

	// Initialize validator
	e.Validator = utils.NewCustomValidator()

	// Setup routes
	routes.SetupRoutes(e)
}
