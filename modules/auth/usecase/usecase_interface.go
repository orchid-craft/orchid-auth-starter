package usecase

import (
	"context"
	"time"

	authReq "github.com/orchid-craft/orchid-auth-starter/modules/auth/delivery/models/request"
	authRes "github.com/orchid-craft/orchid-auth-starter/modules/auth/delivery/models/response"
	"github.com/orchid-craft/orchid-auth-starter/security"
)

// TokenManagerInterface mendefinisikan kontrak pengelolaan token JWT yang dibutuhkan oleh usecase
type TokenManagerInterface interface {
	GenerateAccessToken(userID, email, role string) (string, error)
	GenerateRefreshToken(userID string) (string, string, time.Time, error)
	ValidateRefreshToken(tokenStr string) (*security.RefreshTokenClaims, error)
	HashToken(tokenStr string) string
}

// CSRFManagerInterface mendefinisikan kontrak generasi token CSRF yang dibutuhkan oleh usecase
type CSRFManagerInterface interface {
	GenerateToken() (string, error)
}

type AuthUsecaseInterface interface {
	Register(ctx context.Context, req authReq.RegisterRequest) (*authRes.AuthResponse, error)
	Login(ctx context.Context, req authReq.LoginRequest) (*authRes.AuthResponse, error)
	RefreshToken(ctx context.Context, rawRefreshToken string) (*authRes.RefreshResponse, error)
	Logout(ctx context.Context, rawRefreshToken string) error
	GetCSRFToken(ctx context.Context) (string, error)
}
