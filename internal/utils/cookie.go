package utils

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func SetCookie(ctx echo.Context, name, value, path string, expires time.Time) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		Expires:  expires,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   false, // TODO: set true in production
	}
	ctx.SetCookie(&cookie)
}
