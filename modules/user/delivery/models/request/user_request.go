package request

type CreateUserRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	DisplayName string `json:"display_name"`
}

type UpdateUserRequest struct {
	DisplayName string `json:"display_name"`
}
