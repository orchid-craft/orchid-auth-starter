package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/orchid-craft/orchid-auth-starter/internal/middleware"
	userV1 "github.com/orchid-craft/orchid-auth-starter/modules/user/delivery/api/rest/v1"
	"github.com/orchid-craft/orchid-auth-starter/security"
)

// SetupUserRoutes mengonfigurasi rute modul pengguna dan middleware proteksinya
func SetupUserRoutes(deps *Dependencies) {
	deps.Log.Info("Initialize user routes...")

	jwtMgr := security.NewJWTManager(security.JWTConfig{
		AccessTokenSecret: deps.Cfg.AccessTokenSecret,
		AccessTokenTTL:    deps.Cfg.AccessTokenTTL(),
	})
	csrfMgr := security.NewCSRFManager(deps.Cfg.CSRFSecret)

	h := userV1.NewUserHandler(deps.DI.UserUsecase)

	deps.App.Route("/api", func(r chi.Router) {
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
