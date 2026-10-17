package routes

import (
	httpUtil "github.com/orchid-craft/orchid-auth-starter/internal/bootstrap/server/restful_server/http_util"
)

// SetupCommonRoutes mendaftarkan rute health check dan endpoints umum
func SetupCommonRoutes(deps *Dependencies) {
	deps.App.Get("/health-check", httpUtil.HealthCheckHandler)
	deps.App.Get("/health", httpUtil.HealthCheckHandler)
}
