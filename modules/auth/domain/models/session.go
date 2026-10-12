package models

import "time"

// RefreshTokenSession merepresentasikan status sesi token yang tersimpan di database/in-memory
type RefreshTokenSession struct {
	ID        string    `json:"id"`         // JTI token
	UserID    string    `json:"user_id"`    // ID pengguna
	TokenHash string    `json:"token_hash"` // SHA-256 hash dari refresh token
	Revoked   bool      `json:"revoked"`    // Status pencabutan
	ExpiresAt time.Time `json:"expires_at"` // Waktu kedaluwarsa
	CreatedAt time.Time `json:"created_at"` // Waktu pembuatan
}
