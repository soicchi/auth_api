package controllers

import (
	"log"

	"github.com/soicchi/auth_api/utils"

	"github.com/labstack/echo/v4"
)

type IFUserService interface {
	CreateUser(email, password string) error
}

type UserController struct {
	UserService IFUserService
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserController(us IFUserService) *UserController {
	return &UserController{
		UserService: us,
	}
}

func (c *UserController) SignUp(ctx echo.Context) error {
	var req SignUpRequest
	if err := ctx.Bind(&req); err != nil {
		log.Printf("Failed to bind request: %v", err)
		return utils.BadRequestResponse(ctx, "Invalid request", nil) 
	}

	if err := c.UserService.CreateUser(req.Email, req.Password); err != nil {
		log.Printf("Failed to create user: %v", err)
		return utils.InternalServerErrorResponse(ctx, "Failed to create user", nil)
	}

	return utils.StatusOKResponse(ctx, "Successfully created user", nil)
}
