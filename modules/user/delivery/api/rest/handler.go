package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/yudhiana/web-api/internal/middleware"
	v1 "github.com/yudhiana/web-api/modules/user/delivery/api/rest/v1"
	"github.com/yudhiana/web-api/modules/user/usecase"
)

func NewUserHandler(
	r chi.Router,
	u usecase.UserUsecaseInterface,
	authMiddleware *middleware.AuthMiddleware,
	csrfMiddleware *middleware.CSRFMiddleware,
) {
	h := v1.NewUserHandler(u)

	r.Group(func(userRouter chi.Router) {
		userRouter.Use(authMiddleware.RequireAuth)
		userRouter.Get("/me", h.GetProfile)
		userRouter.Get("/user/profile", h.GetProfile) // alias

		// Admin-only routes
		userRouter.Group(func(adminRouter chi.Router) {
			adminRouter.Use(authMiddleware.RequireRole("admin"))
			adminRouter.Get("/users", h.ListUsers)
		})

		// Mutating routes with CSRF Protection
		userRouter.Group(func(mutating chi.Router) {
			mutating.Use(csrfMiddleware.Protect)
			mutating.Put("/me", h.UpdateProfile)
			mutating.Put("/user/profile", h.UpdateProfile) // alias
		})
	})
}
