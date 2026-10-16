package restfulServer

import (
	"github.com/go-chi/chi/v5"
	"github.com/orchid-craft/orchid-auth-starter/internal/bootstrap/container"
	httpUtil "github.com/orchid-craft/orchid-auth-starter/internal/bootstrap/server/restful_server/http_util"
	authDelivery "github.com/orchid-craft/orchid-auth-starter/modules/auth/delivery/api/rest"
	userDelivery "github.com/orchid-craft/orchid-auth-starter/modules/user/delivery/api/rest"
)

// SetupCommonRoutes configures health checks and general endpoints
func SetupCommonRoutes(app *chi.Mux, container *container.Container) {
	app.Get("/health-check", httpUtil.HealthCheckHandler)
	app.Get("/health", httpUtil.HealthCheckHandler)
}

// SetupAuthRoutes configures the auth module routes
func SetupAuthRoutes(app *chi.Mux, container *container.Container) {
	container.Log.Info("Initialize auth handler...")
	di := container.GetDI()

	app.Route("/api/auth", func(r chi.Router) {
		authDelivery.NewAuthHandler(
			r,
			di.AuthUsecase,
			di.CookieManager,
			di.CSRFManager,
			di.CSRFMiddleware,
			di.RateLimitMiddleware,
		)
	})
}

// SetupUserRoutes configures the user module routes
func SetupUserRoutes(app *chi.Mux, container *container.Container) {
	container.Log.Info("Initialize user handler...")
	di := container.GetDI()

	app.Route("/api", func(r chi.Router) {
		userDelivery.NewUserHandler(
			r,
			di.UserUsecase,
			di.AuthMiddleware,
			di.CSRFMiddleware,
		)
	})
}
