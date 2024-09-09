package utils

import "github.com/go-playground/validator/v10"

func Validate(s interface{}) error {
	valid := validator.New()
	return valid.Struct(s)
}
