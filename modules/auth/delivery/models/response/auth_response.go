package response

import (
	userResponse "github.com/orchid-craft/orchid-auth-starter/modules/user/delivery/models/response"
)

// AccessToken dan RefreshToken sengaja bertag `json:"-"`: token HANYA dikirim lewat cookie HttpOnly
// (di-set oleh handler). Jika token ikut ada di body JSON, script XSS bisa membacanya dan atribut
// HttpOnly + Path isolation pada cookie menjadi tidak berguna. Field tetap ada di struct agar
// handler bisa menyematkannya ke cookie.
type AuthResponse struct {
	User         userResponse.UserResponse `json:"user"`
	AccessToken  string                    `json:"-"`
	RefreshToken string                    `json:"-"`
	CSRFToken    string                    `json:"csrf_token"`
	ExpiresIn    int64                     `json:"expires_in"` // Durasi access token dalam detik
	TokenType    string                    `json:"token_type"`
}

// RefreshResponse: lihat catatan pada AuthResponse mengenai AccessToken dan RefreshToken.
type RefreshResponse struct {
	AccessToken  string `json:"-"`
	RefreshToken string `json:"-"`
	CSRFToken    string `json:"csrf_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}
