package response

import "time"

type UserResponse struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
}
