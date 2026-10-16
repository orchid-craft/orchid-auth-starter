package response

import (
	userResponse "github.com/orchid-craft/orchid-auth-starter/modules/user/delivery/models/response"
)

type AuthResponse struct {
	User         userResponse.UserResponse `json:"user"`
	AccessToken  string                    `json:"access_token,omitempty"`
	RefreshToken string                    `json:"refresh_token,omitempty"`
	CSRFToken    string                    `json:"csrf_token"`
	ExpiresIn    int64                     `json:"expires_in"` // Durasi access token dalam detik
	TokenType    string                    `json:"token_type"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	CSRFToken    string `json:"csrf_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}
