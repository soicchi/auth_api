package utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewResponse(t *testing.T) {
	res := NewResponse(200, "OK", "data")
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "OK", res.Message)
	assert.Equal(t, "data", res.Data)
}

func TestJSONResponse(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/signup", bytes.NewReader([]byte{}))
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)
	res := &Response{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    "TestData",
	}
	err := res.JSONResponse(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "\"message\":\"OK\"")
	assert.Contains(t, rec.Body.String(), "\"data\":\"TestData\"")
}

func TestStatusOKResponse(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/signup", bytes.NewReader([]byte{}))
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)
	err := StatusOKResponse(ctx, "OK", "TestData")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "\"message\":\"OK\"")
	assert.Contains(t, rec.Body.String(), "\"data\":\"TestData\"")
}

func TestBadRequestResponse(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/signup", bytes.NewReader([]byte{}))
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)
	err := BadRequestResponse(ctx, "Bad Request", "TestError")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "\"message\":\"Bad Request\"")
	assert.Contains(t, rec.Body.String(), "\"data\":\"TestError\"")
}

func TestInternalServerErrorResponse(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/sigup", bytes.NewReader([]byte{}))
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)
	err := InternalServerErrorResponse(ctx, "Internal Error", "TestError")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "\"message\":\"Internal Error\"")
	assert.Contains(t, rec.Body.String(), "\"data\":\"TestError\"")
}
