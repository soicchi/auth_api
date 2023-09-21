package routes

import (
	"github.com/soicchi/auth_api/internal/controllers"
	"github.com/soicchi/auth_api/internal/usecase"
	"github.com/soicchi/auth_api/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, cv *utils.CustomValidator) *echo.Echo {
	e := echo.New()

	// Initialize base middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Base routes
	v1 := e.Group("/api/v1")

	// Initialize controllers
	userService := usecase.NewUserService(db, cv)
	userController := controllers.NewUserController(userService)

	// Access token routes
	token := v1.Group("/basic")
	token.POST("/signup", userController.SignUp)

	return e
}
