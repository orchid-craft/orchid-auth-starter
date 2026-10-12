package middleware

import (
	"net/http"

	"github.com/yudhiana/web-api/constants"
	httpUtil "github.com/yudhiana/web-api/internal/bootstrap/server/restful_server/http_util"
	"github.com/yudhiana/web-api/pkg/cookie"
	"github.com/yudhiana/web-api/security"
	bunker "github.com/yudhiana/bunker/errors"
)

type CSRFMiddleware struct {
	csrfManager *security.CSRFManager
}

func NewCSRFMiddleware(csrfManager *security.CSRFManager) *CSRFMiddleware {
	return &CSRFMiddleware{csrfManager: csrfManager}
}

// Protect memvalidasi Double-Submit CSRF token pada metode mutasi HTTP
func (m *CSRFMiddleware) Protect(next http.Handler) http.Handler {
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

		if err := m.csrfManager.ValidateToken(headerToken, cookieVal); err != nil {
			bunkerErr := bunker.New(bunker.StatusForbidden).SetMessage("validasi CSRF gagal: " + err.Error())
			httpUtil.ErrorResponse(w, bunkerErr)
			return
		}

		next.ServeHTTP(w, r)
	})
}
