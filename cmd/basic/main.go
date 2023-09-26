package main

import (
	"log"
	"os"

	"github.com/soicchi/auth_api/internal/models"
	"github.com/soicchi/auth_api/internal/routes"
	"github.com/soicchi/auth_api/internal/utils"
)

func main() {
	// Setup database
	db, err := models.SetupDB()
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	log.Println("Successfully setup database")

	// Setup routes
	e := routes.SetupRoutes(db)
	e.Validator = utils.NewCustomValidator()

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}