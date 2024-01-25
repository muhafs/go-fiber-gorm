package request

import "github.com/go-playground/validator/v10"

// model for receiving requests from client to validate it
type User struct {
	Name     string `json:"name" validate:"required,alpha"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
} // * it's recommended to split User model into (CreateUserSchema) and (UpdateUserSchema).

// model for returned error
type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func (u *User) Validate() []*ErrorResponse {
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
