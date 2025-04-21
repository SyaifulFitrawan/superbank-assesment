package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate

func InitValidator() {
	Validator = validator.New()
	Validator.RegisterValidation("notblank", NotBlank)
}

func NotBlank(fl validator.FieldLevel) bool {
	value := strings.TrimSpace(fl.Field().String())
	return value != ""
}

func ValidateStruct(s interface{}) []string {
	err := Validator.Struct(s)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		field := strings.ToLower(err.Field())
		message := fmt.Sprintf("%s is %s", field, err.Tag())
		messages = append(messages, message)
	}

	return messages
}
