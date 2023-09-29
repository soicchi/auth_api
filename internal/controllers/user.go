package controllers

import (
	"log"

	"github.com/soicchi/auth_api/internal/utils"

	"github.com/labstack/echo/v4"
)

type UserService interface {
	CreateUser(email, password string) error
	CheckSignIn(email, password string) error
}

type UserHandler struct {
	Service UserService
}

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{
		Service: service,
	}
}

func (c *UserHandler) SignUp(ctx echo.Context) error {
	var req SignUpRequest
	if err := ctx.Bind(&req); err != nil {
		log.Printf("Failed to bind request: %v", err)
		return utils.BadRequestResponse(ctx, "Invalid request")
	}

	if err := ctx.Validate(req); err != nil {
		log.Printf("Failed to validate request: %v", err)
		return utils.BadRequestResponse(ctx, "Invalid request")
	}

	if err := c.Service.CreateUser(req.Email, req.Password); err != nil {
		log.Printf("Failed to create user: %v", err)
		return utils.InternalServerErrorResponse(ctx, "Failed to create user")
	}

	return utils.StatusOKResponse(ctx, "Successfully created user", nil)
}

func (c *UserHandler) SignIn(ctx echo.Context) error {
	var req SignInRequest
	if err := ctx.Bind(&req); err != nil {
		log.Printf("Failed to bind request: %v", err)
		return utils.BadRequestResponse(ctx, "Invalid request")
	}

	if err := c.Service.CheckSignIn(req.Email, req.Password); err != nil {
		log.Printf("Failed to sign in: %v", err)
		return utils.BadRequestResponse(ctx, "Invalid email or password")
	}

	return utils.StatusOKResponse(ctx, "Successfully signed in", nil)
}
