package request

import "github.com/go-playground/validator/v10"

// model for returned error
type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func Validate[T any](u T) []*ErrorResponse {
	validate := validator.New()

	var errors []*ErrorResponse

	if err := validate.Struct(u); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse

			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()

			errors = append(errors, &element)
		}
	}

	return errors
}
