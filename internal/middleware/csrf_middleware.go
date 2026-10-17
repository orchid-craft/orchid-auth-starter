package middleware

import (
	"net/http"

	"github.com/orchid-craft/orchid-auth-starter/constants"
	httpUtil "github.com/orchid-craft/orchid-auth-starter/internal/bootstrap/server/restful_server/http_util"
	"github.com/orchid-craft/orchid-auth-starter/pkg/cookie"
	"github.com/orchid-craft/orchid-auth-starter/security"
	bunker "github.com/olympus-forge/bunker/errors"
)

// ProtectCSRF memvalidasi Double-Submit CSRF token pada metode mutasi HTTP
func ProtectCSRF(csrfManager *security.CSRFManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Abaikan safe methods (GET, HEAD, OPTIONS, TRACE)
			switch r.Method {
			case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
				next.ServeHTTP(w, r)
				return
			}

			headerToken := r.Header.Get(constants.HeaderCSRFToken)
			cookieVal := ""
			if c, err := r.Cookie(cookie.CSRFCookieName); err == nil {
				cookieVal = c.Value
			}

			if err := csrfManager.ValidateToken(headerToken, cookieVal); err != nil {
				bunkerErr := httpUtil.NewAppError(bunker.StatusForbidden, "validasi CSRF gagal: "+err.Error())
				httpUtil.ErrorResponse(w, bunkerErr)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

type CSRFMiddleware struct {
	csrfManager *security.CSRFManager
}

func NewCSRFMiddleware(csrfManager *security.CSRFManager) *CSRFMiddleware {
	return &CSRFMiddleware{csrfManager: csrfManager}
}

func (m *CSRFMiddleware) Protect(next http.Handler) http.Handler {
	return ProtectCSRF(m.csrfManager)(next)
}
