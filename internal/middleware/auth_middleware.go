package middleware

import (
	"net/http"
	"strings"

	"github.com/yudhiana/web-api/constants"
	httpUtil "github.com/yudhiana/web-api/internal/bootstrap/server/restful_server/http_util"
	"github.com/yudhiana/web-api/pkg/cookie"
	"github.com/yudhiana/web-api/security"
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
