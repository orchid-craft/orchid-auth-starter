package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	authReq "github.com/orchid-craft/orchid-auth-starter/modules/auth/delivery/models/request"
	authRes "github.com/orchid-craft/orchid-auth-starter/modules/auth/delivery/models/response"
	"github.com/orchid-craft/orchid-auth-starter/modules/auth/domain"
	authModel "github.com/orchid-craft/orchid-auth-starter/modules/auth/domain/models"
	"github.com/orchid-craft/orchid-auth-starter/modules/auth/repository"
	userReq "github.com/orchid-craft/orchid-auth-starter/modules/user/delivery/models/request"
	userDomain "github.com/orchid-craft/orchid-auth-starter/modules/user/domain"
	userUsecase "github.com/orchid-craft/orchid-auth-starter/modules/user/usecase"
	"github.com/orchid-craft/orchid-auth-starter/security"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userUsecase    userUsecase.UserUsecaseInterface
	sessionRepo    repository.SessionRepositoryInterface
	jwtManager     TokenManagerInterface
	csrfManager    CSRFManagerInterface
	accessTokenTTL time.Duration
}

func NewAuthUsecase(
	userUc userUsecase.UserUsecaseInterface,
	sessionRepo repository.SessionRepositoryInterface,
	jwtManager TokenManagerInterface,
	csrfManager CSRFManagerInterface,
	accessTokenTTL time.Duration,
) AuthUsecaseInterface {
	return &authUsecase{
		userUsecase:    userUc,
		sessionRepo:    sessionRepo,
		jwtManager:     jwtManager,
		csrfManager:    csrfManager,
		accessTokenTTL: accessTokenTTL,
	}
}

func (u *authUsecase) Register(ctx context.Context, req authReq.RegisterRequest) (*authRes.AuthResponse, error) {
	// Register user via User Usecase - paksa peran selalu 'user' untuk pendaftaran publik (P0-1)
	_, err := u.userUsecase.Register(ctx, userReq.CreateUserRequest{
		Email:       req.Email,
		Password:    req.Password,
		Role:        "user",
		DisplayName: req.DisplayName,
	})
	if err != nil {
		return nil, err
	}

	// Login immediately after successful registration
	return u.Login(ctx, authReq.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
}

func (u *authUsecase) Login(ctx context.Context, req authReq.LoginRequest) (*authRes.AuthResponse, error) {
	user, err := u.userUsecase.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, userDomain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := u.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat access token: %w", err)
	}

	refreshToken, tokenID, expiresAt, err := u.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat refresh token: %w", err)
	}

	session := &authModel.RefreshTokenSession{
		ID:        tokenID,
		UserID:    user.ID,
		TokenHash: u.jwtManager.HashToken(refreshToken),
		Revoked:   false,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	if err := u.sessionRepo.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("gagal menyimpan sesi refresh token: %w", err)
	}

	csrfToken, err := u.csrfManager.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("gagal membuat csrf token: %w", err)
	}

	userRes, err := u.userUsecase.GetByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &authRes.AuthResponse{
		User:         *userRes,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		CSRFToken:    csrfToken,
		ExpiresIn:    int64(u.accessTokenTTL.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (u *authUsecase) RefreshToken(ctx context.Context, rawRefreshToken string) (*authRes.RefreshResponse, error) {
	if rawRefreshToken == "" {
		return nil, domain.ErrMissingRefreshToken
	}

	claims, err := u.jwtManager.ValidateRefreshToken(rawRefreshToken)
	if err != nil {
		if errors.Is(err, security.ErrTokenExpired) {
			return nil, domain.ErrTokenExpired
		}
		return nil, domain.ErrTokenInvalid
	}

	session, err := u.sessionRepo.FindByID(ctx, claims.TokenID)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, domain.ErrTokenExpired
	}

	// Deteksi reuse attempt: Jika token sudah dicabut, cabut seluruh keluarga token user
	if session.Revoked {
		_ = u.sessionRepo.RevokeAllForUser(ctx, claims.UserID)
		return nil, domain.ErrTokenReused
	}

	if session.TokenHash != u.jwtManager.HashToken(rawRefreshToken) {
		return nil, domain.ErrTokenInvalid
	}

	// Revoke token lama (Refresh Token Rotation)
	if err := u.sessionRepo.Revoke(ctx, claims.TokenID); err != nil {
		return nil, fmt.Errorf("gagal mencabut token lama: %w", err)
	}

	user, err := u.userUsecase.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, userDomain.ErrUserNotFound
	}

	newAccessToken, err := u.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	newRefreshToken, newTokenID, expiresAt, err := u.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	newSession := &authModel.RefreshTokenSession{
		ID:        newTokenID,
		UserID:    user.ID,
		TokenHash: u.jwtManager.HashToken(newRefreshToken),
		Revoked:   false,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	if err := u.sessionRepo.Save(ctx, newSession); err != nil {
		return nil, fmt.Errorf("gagal menyimpan sesi refresh token baru: %w", err)
	}

	newCSRFToken, err := u.csrfManager.GenerateToken()
	if err != nil {
		return nil, err
	}

	return &authRes.RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		CSRFToken:    newCSRFToken,
		ExpiresIn:    int64(u.accessTokenTTL.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (u *authUsecase) Logout(ctx context.Context, rawRefreshToken string) error {
	if rawRefreshToken == "" {
		return nil
	}

	claims, err := u.jwtManager.ValidateRefreshToken(rawRefreshToken)
	if err != nil {
		return nil
	}

	return u.sessionRepo.Revoke(ctx, claims.TokenID)
}

func (u *authUsecase) GetCSRFToken(ctx context.Context) (string, error) {
	return u.csrfManager.GenerateToken()
}
