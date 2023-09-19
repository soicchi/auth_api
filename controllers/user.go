package controllers

import (
	"log"
	"net/http"

	"github.com/soicchi/auth_api/usecase"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	UserService *usecase.UserService
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserController(us *usecase.UserService) *UserController {
	return &UserController{
		UserService: us,
	}
}

func (c *UserController) SignUp(ctx echo.Context) error {
	var req SignUpRequest
	if err := ctx.Bind(&req); err != nil {
		log.Printf("Failed to bind request: %v", err)
		return ctx.JSON(http.StatusBadRequest, "Invalid request")
	}

	if err := c.UserService.CreateUser(req.Email, req.Password); err != nil {
		log.Printf("Failed to create user: %v", err)
		return ctx.JSON(http.StatusInternalServerError, "Failed to create user")
	}

	return ctx.JSON(http.StatusOK, "Successfully created user")
}
