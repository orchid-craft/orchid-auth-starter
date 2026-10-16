package restfulServer

import (
	"github.com/go-chi/chi/v5"
	httpUtil "github.com/orchid-craft/orchid-auth-starter/internal/bootstrap/server/restful_server/http_util"
	"github.com/orchid-craft/orchid-auth-starter/internal/middleware"
	authV1 "github.com/orchid-craft/orchid-auth-starter/modules/auth/delivery/api/rest/v1"
	userV1 "github.com/orchid-craft/orchid-auth-starter/modules/user/delivery/api/rest/v1"
	"github.com/orchid-craft/orchid-auth-starter/pkg/cookie"
	"github.com/orchid-craft/orchid-auth-starter/security"
)

// SetupCommonRoutes configures health checks and general endpoints
func SetupCommonRoutes(s *Server) {
	s.app.Get("/health-check", httpUtil.HealthCheckHandler)
	s.app.Get("/health", httpUtil.HealthCheckHandler)
}

// SetupAuthRoutes configures the auth module routes and mounts middlewares directly on the router
func SetupAuthRoutes(s *Server) {
	s.log.Info("Initialize auth routes...")

	cookieMgr := cookie.NewManager(cookie.Config{
		Secure:          s.cfg.CookieSecure,
		Domain:          s.cfg.CookieDomain,
		AccessTokenTTL:  s.cfg.AccessTokenTTL(),
		RefreshTokenTTL: s.cfg.RefreshTokenTTL(),
	})
	csrfMgr := security.NewCSRFManager(s.cfg.CSRFSecret)

	h := authV1.NewAuthHandler(s.di.AuthUsecase, cookieMgr)

	s.app.Route("/api/auth", func(r chi.Router) {
		// Public auth endpoints
		r.Get("/csrf-token", h.CSRFToken)
		r.Post("/register", h.Register)

		// Login dengan proteksi rate limit anti brute-force
		if s.rateLimiter != nil {
			r.With(s.rateLimiter.Limit).Post("/login", h.Login)
		} else {
			r.Post("/login", h.Login)
		}

		// Endpoints yang dilindungi validasi CSRF
		r.Group(func(csrfProtected chi.Router) {
			csrfProtected.Use(middleware.ProtectCSRF(csrfMgr))

			// Refresh token dengan rate limit
			if s.rateLimiter != nil {
				csrfProtected.With(s.rateLimiter.Limit).Post("/refresh", h.RefreshToken)
			} else {
				csrfProtected.Post("/refresh", h.RefreshToken)
			}

			csrfProtected.Post("/logout", h.Logout)
		})
	})
}

// SetupUserRoutes configures the user module routes and mounts middlewares directly on the router
func SetupUserRoutes(s *Server) {
	s.log.Info("Initialize user routes...")

	jwtMgr := security.NewJWTManager(security.JWTConfig{
		AccessTokenSecret: s.cfg.AccessTokenSecret,
		AccessTokenTTL:    s.cfg.AccessTokenTTL(),
	})
	csrfMgr := security.NewCSRFManager(s.cfg.CSRFSecret)

	h := userV1.NewUserHandler(s.di.UserUsecase)

	s.app.Route("/api", func(r chi.Router) {
		// Endpoints yang dilindungi Access Token JWT
		r.Group(func(authProtected chi.Router) {
			authProtected.Use(middleware.RequireAuth(jwtMgr))

			authProtected.Get("/me", h.GetProfile)
			authProtected.Get("/user/profile", h.GetProfile) // alias kompatibel

			// Admin-only route
			authProtected.With(middleware.RequireRole("admin")).Get("/users", h.ListUsers)

			// Mutating profile route dilindungi CSRF
			authProtected.Group(func(mutating chi.Router) {
				mutating.Use(middleware.ProtectCSRF(csrfMgr))
				mutating.Put("/me", h.UpdateProfile)
				mutating.Put("/user/profile", h.UpdateProfile) // alias kompatibel
			})
		})
	})
}
