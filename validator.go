package configs

import (
	"github.com/go-playground/validator"
)

type Validator interface {
	Validate(config interface{}) error
}

type StructValidator struct {
	validator *validator.Validate
}

func NewValidator(validator *validator.Validate) Validator {
	return &StructValidator{
		validator: validator,
	}
}

func (that *StructValidator) Validate(config interface{}) error {
	return that.validator.Struct(config)
}
