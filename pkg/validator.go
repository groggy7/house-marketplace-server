package pkg

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s any) []string {
	err := validate.Struct(s)
	if err != nil {
		var errorMessages []string
		for _, fieldError := range err.(validator.ValidationErrors) {
			switch fieldError.Tag() {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("%s field is required.", fieldError.Field()))
			case "min":
				errorMessages = append(errorMessages, fmt.Sprintf("%s field must be at least %s characters long.", fieldError.Field(), fieldError.Param()))
			case "max":
				errorMessages = append(errorMessages, fmt.Sprintf("%s field must be at most %s characters long.", fieldError.Field(), fieldError.Param()))
			case "oneof":
				errorMessages = append(errorMessages, fmt.Sprintf("%s field must be one of: %s", fieldError.Field(), fieldError.Param()))
			}
		}
		return errorMessages
	}
	return nil
}
