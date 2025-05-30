package utils

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateRegisterInput(req interface{}) error {
	if err := validate.Struct(req); err != nil {
		return err
	}
	return nil
}

func ValidateUpdateProfileInput(req interface{}) error {
	if err := validate.Struct(req); err != nil {
		return err
	}
	return nil
}

func ParseValidationError(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		return errs.Error()
	}
	return err.Error()
}
