package utils

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{
		Validator: validator.New(),
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return fmt.Errorf("error validating struct %v", err)
	}
	return nil
}

func ValidateENVVars(vars []string) error {
	for _, v := range vars {
		if os.Getenv(v) == "" {
			return fmt.Errorf("%s environment variable not set", v)
		}
	}
	return nil
}

func NewENVVars() []string {
	return []string{
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_NAME",
		"DB_PASSWORD",
		"DB_SSL_MODE",
		"API_PORT",
		"BASIC_AUTH_USERNAME",
		"BASIC_AUTH_PASSWORD",
		"API_KEY",
	}
}
