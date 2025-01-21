package configValidator

import (
	"github.com/go-playground/validator"
)

type StructValidator struct {
	validator *validator.Validate
}

func New(validator *validator.Validate) Validator {
	return &StructValidator{
		validator: validator,
	}
}

func (that *StructValidator) Validate(config interface{}) error {
	return that.validator.Struct(config)
}
