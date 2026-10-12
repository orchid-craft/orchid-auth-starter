package request

type CreateUserRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8,max=100"`
	Role        string `json:"role" validate:"omitempty,oneof=user admin"`
	DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
}

type UpdateUserRequest struct {
	DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
}
