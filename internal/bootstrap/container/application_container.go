package container

import (
	"context"

	"github.com/olympus-forge/logos"
	"github.com/orchid-craft/orchid-auth-starter/config"
	authRepo "github.com/orchid-craft/orchid-auth-starter/modules/auth/repository"
	authUsecase "github.com/orchid-craft/orchid-auth-starter/modules/auth/usecase"
	userRepo "github.com/orchid-craft/orchid-auth-starter/modules/user/repository"
	userUsecase "github.com/orchid-craft/orchid-auth-starter/modules/user/usecase"
	userReq "github.com/orchid-craft/orchid-auth-starter/modules/user/delivery/models/request"
	"github.com/orchid-craft/orchid-auth-starter/security"
)

type DirectInjection struct {
	UserRepo    userRepo.UserRepositoryInterface
	SessionRepo authRepo.SessionRepositoryInterface
	UserUsecase userUsecase.UserUsecaseInterface
	AuthUsecase authUsecase.AuthUsecaseInterface
}

func NewApplicationContainer(cfg *config.LocalConfig) (*DirectInjection, error) {
	logger := logos.NewLogger()
	logger.Info("Initializing dependency injection container...")

	// 1. Inisialisasi Security Managers (JWT & CSRF)
	jwtConfig := security.JWTConfig{
		AccessTokenSecret:  cfg.AccessTokenSecret,
		RefreshTokenSecret: cfg.RefreshTokenSecret,
		AccessTokenTTL:     cfg.AccessTokenTTL(),
		RefreshTokenTTL:    cfg.RefreshTokenTTL(),
	}
	jwtManager := security.NewJWTManager(jwtConfig)
	csrfManager := security.NewCSRFManager(cfg.CSRFSecret)

	// 2. Inisialisasi Repositories (In-Memory atau Database)
	uRepo := userRepo.NewUserRepository()
	sRepo := authRepo.NewSessionRepository()

	// 3. Inisialisasi Usecases
	uUsecase := userUsecase.NewUserUsecase(uRepo)
	aUsecase := authUsecase.NewAuthUsecase(uUsecase, sRepo, jwtManager, csrfManager, cfg.AccessTokenTTL())

	// Seed akun demo (alice@example.com = admin, password publik) HANYA untuk development/test.
	// Di production, akun admin dengan password yang diketahui siapa pun berarti pengambilalihan penuh.
	if cfg.IsDevelopment() {
		seedDefaultUsers(uUsecase)
	} else {
		logger.Info("Seed akun default dilewati (APP_ENV bukan development/test)")
	}

	di := &DirectInjection{
		UserRepo:    uRepo,
		SessionRepo: sRepo,
		UserUsecase: uUsecase,
		AuthUsecase: aUsecase,
	}

	logger.Info("Dependency injection container initialized successfully")
	return di, nil
}

func seedDefaultUsers(u userUsecase.UserUsecaseInterface) {
	ctx := context.Background()
	_, _ = u.Register(ctx, userReq.CreateUserRequest{
		Email:       "alice@example.com",
		Password:    "password123",
		DisplayName: "Alice Admin",
		Role:        "admin",
	})
	_, _ = u.Register(ctx, userReq.CreateUserRequest{
		Email:       "bob@example.com",
		Password:    "password123",
		DisplayName: "Bob Member",
		Role:        "user",
	})
}

func (di *DirectInjection) Close() error {
	logger := logos.NewLogger()
	logger.Info("🔴Closing dependency injection resources...")
	return nil
}
