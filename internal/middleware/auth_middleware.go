package middleware

import (
	"net/http"
	"strings"

	"github.com/orchid-craft/orchid-auth-starter/constants"
	httpUtil "github.com/orchid-craft/orchid-auth-starter/internal/bootstrap/server/restful_server/http_util"
	"github.com/orchid-craft/orchid-auth-starter/pkg/cookie"
	"github.com/orchid-craft/orchid-auth-starter/security"
	bunker "github.com/yudhiana/bunker/errors"
)

type AuthMiddleware struct {
	jwtManager *security.JWTManager
}

func NewAuthMiddleware(jwtManager *security.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{jwtManager: jwtManager}
}

// RequireAuth memvalidasi JWT dari Authorization Header atau HttpOnly Cookie
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := ""

		authHeader := r.Header.Get(constants.HeaderAuthorization)
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			if c, err := r.Cookie(cookie.AccessTokenCookieName); err == nil && c.Value != "" {
				tokenStr = c.Value
			}
		}

		if tokenStr == "" {
			err := bunker.New(bunker.StatusUnAuthorized).SetMessage("autentikasi dibutuhkan: token tidak ditemukan")
			httpUtil.ErrorResponse(w, err)
			return
		}

		claims, err := m.jwtManager.ValidateAccessToken(tokenStr)
		if err != nil {
			status := bunker.StatusUnAuthorized
			msg := "token tidak valid atau telah kedaluwarsa"
			if err == security.ErrTokenExpired {
				msg = "token telah kedaluwarsa"
			}
			errBunker := bunker.New(status).SetMessage(msg)
			httpUtil.ErrorResponse(w, errBunker)
			return
		}

		ctx := security.SetClaimsContext(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole membatasi akses hanya untuk pengguna dengan peran tertentu (RBAC)
func (m *AuthMiddleware) RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := security.GetClaimsContext(r.Context())
			if !ok || claims == nil {
				err := bunker.New(bunker.StatusUnAuthorized).SetMessage("autentikasi dibutuhkan: klaim pengguna tidak ditemukan")
				httpUtil.ErrorResponse(w, err)
				return
			}

			isAllowed := false
			for _, role := range allowedRoles {
				if claims.Role == role {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				err := bunker.New(bunker.StatusForbidden).SetMessage("akses ditolak: peran pengguna tidak memiliki izin yang cukup")
				httpUtil.ErrorResponse(w, err)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
