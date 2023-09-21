package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Code    int
	Message string
	Data    interface{}
}

func NewResponse(code int, message string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func (res *Response) JSONResponse(ctx echo.Context) error {
	return ctx.JSON(res.Code, map[string]interface{}{
		"message": res.Message,
		"data":    res.Data,
	})
}

func StatusOKResponse(ctx echo.Context, message string, data interface{}) error {
	res := NewResponse(http.StatusOK, message, data)
	return res.JSONResponse(ctx)
}

func BadRequestResponse(ctx echo.Context, message string, data interface{}) error {
	res := NewResponse(http.StatusBadRequest, message, data)
	return res.JSONResponse(ctx)
}

func InternalServerErrorResponse(ctx echo.Context, message string, data interface{}) error {
	res := NewResponse(http.StatusInternalServerError, message, data)
	return res.JSONResponse(ctx)
}
