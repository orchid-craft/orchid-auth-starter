package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/orchid-craft/orchid-auth-starter/internal/middleware"
	authV1 "github.com/orchid-craft/orchid-auth-starter/modules/auth/delivery/api/rest/v1"
	"github.com/orchid-craft/orchid-auth-starter/pkg/cookie"
	"github.com/orchid-craft/orchid-auth-starter/security"
)

// SetupAuthRoutes mengonfigurasi rute autentikasi dan middleware proteksinya
func SetupAuthRoutes(deps *Dependencies) {
	deps.Log.Info("Initialize auth routes...")

	cookieMgr := cookie.NewManager(cookie.Config{
		Secure:          deps.Cfg.CookieSecure,
		Domain:          deps.Cfg.CookieDomain,
		AccessTokenTTL:  deps.Cfg.AccessTokenTTL(),
		RefreshTokenTTL: deps.Cfg.RefreshTokenTTL(),
	})
	csrfMgr := security.NewCSRFManager(deps.Cfg.CSRFSecret)

	h := authV1.NewAuthHandler(deps.DI.AuthUsecase, cookieMgr)

	deps.App.Route("/api/auth", func(r chi.Router) {
		// Public auth endpoints
		r.Get("/csrf-token", h.CSRFToken)
		r.Post("/register", h.Register)

		// Login dengan proteksi rate limit anti brute-force
		if deps.RateLimiter != nil {
			r.With(deps.RateLimiter.Limit).Post("/login", h.Login)
		} else {
			r.Post("/login", h.Login)
		}

		// Endpoints yang dilindungi validasi CSRF
		r.Group(func(csrfProtected chi.Router) {
			csrfProtected.Use(middleware.ProtectCSRF(csrfMgr))

			// Refresh token dengan rate limit
			if deps.RateLimiter != nil {
				csrfProtected.With(deps.RateLimiter.Limit).Post("/refresh", h.RefreshToken)
			} else {
				csrfProtected.Post("/refresh", h.RefreshToken)
			}

			csrfProtected.Post("/logout", h.Logout)
		})
	})
}
