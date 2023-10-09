package routes

import (
	"github.com/soicchi/auth_api/internal/controllers"
	"github.com/soicchi/auth_api/internal/middleware"
	"github.com/soicchi/auth_api/internal/models"
	"github.com/soicchi/auth_api/internal/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func setupV1Routes(v1 *echo.Group, db *gorm.DB) {
	// Initialize user handler
	userRepo := models.NewUserPostgresRepository(db)
	tokenRepo := models.NewRefreshTokenPostgresRepository(db)
	userService := usecase.NewUserServiceImpl(userRepo, tokenRepo)
	userHandler := controllers.NewUserHandler(userService)

	// Basic Auth
	basic := v1.Group("/basic")
	basic.Use(middleware.BasicAuth)
	basic.POST("/users", userHandler.ListUsers)

	refreshTokenRepo := models.NewRefreshTokenPostgresRepository(db)
	refreshTokenService := usecase.NewRefreshTokenServiceImpl(refreshTokenRepo)
	refreshTokenHandler := controllers.NewRefreshTokenHandler(refreshTokenService)
	// Key Auth
	key := v1.Group("/key")
	key.Use(middleware.KeyAuth)
	key.POST("/signup", userHandler.SignUp)
	key.POST("/signin", userHandler.SignIn)
	key.POST("/refresh", refreshTokenHandler.PostRefreshToken)

	// JWT Auth
	jwt := v1.Group("/jwt")
	jwt.Use(middleware.JWTAuth)
	jwt.GET("/users", userHandler.ListUsers)
}
