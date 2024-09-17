package common

import "github.com/go-playground/validator/v10"

func ValidationErrors(err error) []string {
	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}
	return errors
}
