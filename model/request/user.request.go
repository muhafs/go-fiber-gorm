package request

// model for receiving requests from client to validate it
type User struct {
	Name     string `json:"name" validate:"required,alpha"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
} // * it's recommended to split User model into (CreateUserSchema) and (UpdateUserSchema).
