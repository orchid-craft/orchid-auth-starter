package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("pengguna tidak ditemukan")
	ErrUserAlreadyExists = errors.New("email sudah terdaftar")
	ErrInvalidPassword   = errors.New("kata sandi tidak cocok")
)
