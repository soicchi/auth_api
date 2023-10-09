package controllers

import (
	"log"

	"github.com/soicchi/auth_api/internal/utils"

	"github.com/labstack/echo/v4"
)

type RefreshTokenHandler struct {
	Service RefreshTokenService
}

type RefreshTokenService interface {
	RefreshAccessToken(token string) (string, error)
}

type refreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func NewRefreshTokenHandler(service RefreshTokenService) RefreshTokenHandler {
	return RefreshTokenHandler{
		Service: service,
	}
}

func newRefreshTokenResponse(accessToken string) refreshTokenResponse {
	return refreshTokenResponse{
		AccessToken: accessToken,
	}
}

func (h *RefreshTokenHandler) PostRefreshToken(ctx echo.Context) error {
	token, err := ctx.Cookie("refresh_token")
	if err != nil {
		log.Printf("failed to get cookie: %v", err)
		return utils.BadRequestResponse(ctx, "bad request")
	}

	accessToken, err := h.Service.RefreshAccessToken(token.Value)
	if err != nil {
		log.Printf("failed to refresh access token: %v", err)
		return utils.BadRequestResponse(ctx, "bad request")
	}

	response := newRefreshTokenResponse(accessToken)
	return utils.StatusOKResponse(ctx, "New access token has been issued", response)
}
