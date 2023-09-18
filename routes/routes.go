package routes

import (
	"github.com/soicchi/auth_api/controllers"
	"github.com/soicchi/auth_api/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	// Initialize base middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Base routes
	v1 := e.Group("/api/v1")

	// Initialize controllers
	us := usecase.NewUserService(db, e.Validator)
	uc := controllers.NewUserController(us)

	// Access token routes
	token := v1.Group("/token")
	token.POST("/signup", uc.SignUpHandler)
}
