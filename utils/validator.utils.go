package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// NewValidator func for create a new validator for model fields.
func NewValidator() *validator.Validate {
	// Create a new validator for a Book model.
	validate := validator.New()

	// Custom validation for uuid.UUID fields.
	_ = validate.RegisterValidation("myuuid", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		if _, err := uuid.Parse(field); err != nil {
			return false
		}
		return true
	})

	return validate
}

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value any    `json:"value,omitempty"`
	Error string `json:"error"`
}

// ValidatorErrors func for show validation errors for each invalid fields.
func ValidatorErrors(err error) []*ErrorResponse {
	// define errors struct
	var fieldErrors []*ErrorResponse

	// Make error message for each invalid field.
	for _, err := range err.(validator.ValidationErrors) {
		// Define error fields.
		var element ErrorResponse

		element.Field = err.Field()
		element.Tag = err.Tag()
		element.Value = err.Value()
		element.Error = err.Error()

		fieldErrors = append(fieldErrors, &element)
	}

	return fieldErrors
}
