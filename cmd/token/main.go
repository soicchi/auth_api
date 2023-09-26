package main

import (
	"log"
	"os"

	"github.com/soicchi/auth_api/internal/models"
	"github.com/soicchi/auth_api/internal/routes"
	"github.com/soicchi/auth_api/internal/utils"
)

func main() {
	// Create database connection
	db, err := models.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	log.Println("Successfully connected to database")

	// Migrate database
	models.Migrate(db)

	log.Println("Successfully migrated database")

	// Setup routes
	e := routes.SetupRoutes(db)
	e.Validator = utils.NewCustomValidator()

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
