package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitializeMiddleware(e *echo.Echo) {
	// Initialize base middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
}
