package container

import (
	"context"

	"github.com/olympus-forge/logos"
	"github.com/yudhiana/web-api/config"
	"github.com/yudhiana/web-api/internal/middleware"
	authRepo "github.com/yudhiana/web-api/modules/auth/repository"
	authUsecase "github.com/yudhiana/web-api/modules/auth/usecase"
	userRepo "github.com/yudhiana/web-api/modules/user/repository"
	userUsecase "github.com/yudhiana/web-api/modules/user/usecase"
	userReq "github.com/yudhiana/web-api/modules/user/delivery/models/request"
	"github.com/yudhiana/web-api/pkg/cookie"
	"github.com/yudhiana/web-api/security"
)

type DirectInjection struct {
	UserRepo       userRepo.UserRepositoryInterface
	SessionRepo    authRepo.SessionRepositoryInterface
	UserUsecase    userUsecase.UserUsecaseInterface
	AuthUsecase    authUsecase.AuthUsecaseInterface
	JWTManager     *security.JWTManager
	CSRFManager    *security.CSRFManager
	CookieManager  *cookie.Manager
	AuthMiddleware *middleware.AuthMiddleware
	CSRFMiddleware *middleware.CSRFMiddleware
}

func NewApplicationContainer(cfg *config.LocalConfig) (*DirectInjection, error) {
	logger := logos.NewLogger()
	logger.Info("Initializing dependency injection container...")

	// 1. Inisialisasi Security & Cookie Managers
	jwtConfig := security.JWTConfig{
		AccessTokenSecret:  cfg.AccessTokenSecret,
		RefreshTokenSecret: cfg.RefreshTokenSecret,
		AccessTokenTTL:     cfg.AccessTokenTTL(),
		RefreshTokenTTL:    cfg.RefreshTokenTTL(),
	}
	jwtManager := security.NewJWTManager(jwtConfig)
	csrfManager := security.NewCSRFManager(cfg.CSRFSecret)

	cookieConfig := cookie.Config{
		Secure:          cfg.CookieSecure,
		Domain:          cfg.CookieDomain,
		AccessTokenTTL:  cfg.AccessTokenTTL(),
		RefreshTokenTTL: cfg.RefreshTokenTTL(),
	}
	cookieManager := cookie.NewManager(cookieConfig)

	// 2. Inisialisasi Middlewares
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)
	csrfMiddleware := middleware.NewCSRFMiddleware(csrfManager)

	// 3. Inisialisasi Repositories (In-Memory atau Database)
	uRepo := userRepo.NewUserRepository()
	sRepo := authRepo.NewSessionRepository()

	// 4. Inisialisasi Usecases
	uUsecase := userUsecase.NewUserUsecase(uRepo)
	aUsecase := authUsecase.NewAuthUsecase(uUsecase, sRepo, jwtManager, csrfManager, cfg.AccessTokenTTL())

	// Seed akun default untuk testing/development
	seedDefaultUsers(uUsecase)

	di := &DirectInjection{
		UserRepo:       uRepo,
		SessionRepo:    sRepo,
		UserUsecase:    uUsecase,
		AuthUsecase:    aUsecase,
		JWTManager:     jwtManager,
		CSRFManager:    csrfManager,
		CookieManager:  cookieManager,
		AuthMiddleware: authMiddleware,
		CSRFMiddleware: csrfMiddleware,
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
