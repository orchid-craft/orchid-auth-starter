package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/olympus-forge/logos"
	"github.com/orchid-craft/orchid-auth-starter/config"
	"github.com/orchid-craft/orchid-auth-starter/internal/bootstrap/container"
	"github.com/orchid-craft/orchid-auth-starter/internal/middleware"
)

// Dependencies berisi dependensi yang dibutuhkan oleh router modul
type Dependencies struct {
	App         *chi.Mux
	Cfg         *config.LocalConfig
	DI          *container.DirectInjection
	Log         *logos.LogEntry
	RateLimiter *middleware.RateLimiterMiddleware
}

// RouteSetup mendefinisikan signature fungsi untuk mendaftarkan rute modul
type RouteSetup func(deps *Dependencies)

// SetupAllRoutes mengonfigurasi seluruh rute aplikasi secara terpusat
func SetupAllRoutes(deps *Dependencies, routeSetups ...RouteSetup) {
	deps.Log.Info("Setting up all application routes...")

	if len(routeSetups) == 0 {
		routeSetups = []RouteSetup{
			SetupCommonRoutes,
			SetupAuthRoutes,
			SetupUserRoutes,
		}
	}

	for _, setupFunc := range routeSetups {
		setupFunc(deps)
	}

	deps.Log.Info("All application routes configured successfully")
}
