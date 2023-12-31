package routes

import (
	"github.com/soicchi/auth_api/internal/middleware"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *echo.Echo {
	e := echo.New()

	// Initialize base middleware
	middleware.InitializeMiddleware(e)

	// Setup v1 routes
	v1 := e.Group("/api/v1")
	setupV1Routes(v1, db)

	return e
}
