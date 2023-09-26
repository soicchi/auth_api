package routes

import (
	"github.com/soicchi/auth_api/internal/controllers"
	"github.com/soicchi/auth_api/internal/models"
	"github.com/soicchi/auth_api/internal/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func setupV1Routes(v1 *echo.Group, db *gorm.DB) {
	// Initialize user handler
	userRepo := models.NewUserPostgresRepository(db)
	userService := usecase.NewUserServiceImpl(userRepo)
	userHandler := controllers.NewUserHandler(userService)

	basic := v1.Group("/basic")
	basic.POST("/signup", userHandler.SignUp)
}
