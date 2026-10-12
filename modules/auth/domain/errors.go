package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("email atau kata sandi tidak valid")
	ErrTokenExpired        = errors.New("token telah kedaluwarsa")
	ErrTokenInvalid        = errors.New("token tidak valid")
	ErrTokenRevoked        = errors.New("token telah dicabut")
	ErrTokenReused         = errors.New("peringatan keamanan: refresh token terdeteksi digunakan ulang (reuse attempt)")
	ErrMissingRefreshToken = errors.New("refresh token tidak ditemukan")
)
