package usecase

import (
	"context"

	authReq "github.com/orchid-craft/orchid-auth-starter/modules/auth/delivery/models/request"
	authRes "github.com/orchid-craft/orchid-auth-starter/modules/auth/delivery/models/response"
)

type AuthUsecaseInterface interface {
	Register(ctx context.Context, req authReq.RegisterRequest) (*authRes.AuthResponse, error)
	Login(ctx context.Context, req authReq.LoginRequest) (*authRes.AuthResponse, error)
	RefreshToken(ctx context.Context, rawRefreshToken string) (*authRes.RefreshResponse, error)
	Logout(ctx context.Context, rawRefreshToken string) error
}
