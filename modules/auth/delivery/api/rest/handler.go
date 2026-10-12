package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/yudhiana/web-api/internal/middleware"
	v1 "github.com/yudhiana/web-api/modules/auth/delivery/api/rest/v1"
	"github.com/yudhiana/web-api/modules/auth/usecase"
	"github.com/yudhiana/web-api/pkg/cookie"
	"github.com/yudhiana/web-api/security"
)

func NewAuthHandler(
	r chi.Router,
	u usecase.AuthUsecaseInterface,
	cm *cookie.Manager,
	csrfManager *security.CSRFManager,
	csrfMiddleware *middleware.CSRFMiddleware,
	rateLimitMiddleware *middleware.RateLimiterMiddleware,
) {
	h := v1.NewAuthHandler(u, cm, csrfManager)

	r.Group(func(authRouter chi.Router) {
		// Public auth endpoints
		authRouter.Get("/csrf-token", h.CSRFToken)
		authRouter.Post("/register", h.Register)

		// Login dengan perlindungan rate-limit anti brute-force (P1-3)
		if rateLimitMiddleware != nil {
			authRouter.With(rateLimitMiddleware.Limit).Post("/login", h.Login)
		} else {
			authRouter.Post("/login", h.Login)
		}

		// Endpoints protected by CSRF validation
		authRouter.Group(func(protected chi.Router) {
			protected.Use(csrfMiddleware.Protect)

			// Refresh dengan perlindungan rate-limit (P1-3)
			if rateLimitMiddleware != nil {
				protected.With(rateLimitMiddleware.Limit).Post("/refresh", h.RefreshToken)
			} else {
				protected.Post("/refresh", h.RefreshToken)
			}

			protected.Post("/logout", h.Logout)
		})
	})
}
