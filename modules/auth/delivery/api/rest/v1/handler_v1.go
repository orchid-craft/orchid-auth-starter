package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	bunker "github.com/yudhiana/bunker/errors"
	httpUtil "github.com/orchid-craft/orchid-auth-starter/internal/bootstrap/server/restful_server/http_util"
	authReq "github.com/orchid-craft/orchid-auth-starter/modules/auth/delivery/models/request"
	"github.com/orchid-craft/orchid-auth-starter/modules/auth/domain"
	"github.com/orchid-craft/orchid-auth-starter/modules/auth/usecase"
	userDomain "github.com/orchid-craft/orchid-auth-starter/modules/user/domain"
	"github.com/orchid-craft/orchid-auth-starter/pkg/cookie"
	"github.com/orchid-craft/orchid-auth-starter/pkg/validator"
	"github.com/orchid-craft/orchid-auth-starter/security"
)

type authHandler struct {
	usecase       usecase.AuthUsecaseInterface
	cookieManager *cookie.Manager
	csrfManager   *security.CSRFManager
}

func NewAuthHandler(u usecase.AuthUsecaseInterface, cm *cookie.Manager, csrf *security.CSRFManager) *authHandler {
	return &authHandler{
		usecase:       u,
		cookieManager: cm,
		csrfManager:   csrf,
	}
}

// Register menangani pendaftaran pengguna baru
func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req authReq.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errBunker := bunker.New(bunker.StatusBadRequest).SetMessage("format request payload tidak valid")
		httpUtil.ErrorResponse(w, errBunker)
		return
	}

	if err := validator.ValidateStruct(&req); err != nil {
		httpUtil.ErrorResponse(w, bunker.New(bunker.StatusBadRequest).SetMessage(err.Error()))
		return
	}

	res, err := h.usecase.Register(r.Context(), req)
	if err != nil {
		status := bunker.StatusBadRequest
		if errors.Is(err, userDomain.ErrUserAlreadyExists) {
			status = bunker.StatusConflict
		}
		httpUtil.ErrorResponse(w, bunker.New(status).SetMessage(err.Error()))
		return
	}

	h.cookieManager.SetAccessTokenCookie(w, res.AccessToken)
	h.cookieManager.SetRefreshTokenCookie(w, res.RefreshToken)
	h.cookieManager.SetCSRFCookie(w, res.CSRFToken)

	httpUtil.SuccessResponse(w, res)
}

// Login memvalidasi login pengguna dan menyematkan cookies
func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authReq.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errBunker := bunker.New(bunker.StatusBadRequest).SetMessage("format request payload tidak valid")
		httpUtil.ErrorResponse(w, errBunker)
		return
	}

	if err := validator.ValidateStruct(&req); err != nil {
		httpUtil.ErrorResponse(w, bunker.New(bunker.StatusBadRequest).SetMessage(err.Error()))
		return
	}

	res, err := h.usecase.Login(r.Context(), req)
	if err != nil {
		status := bunker.StatusBadRequest
		if errors.Is(err, domain.ErrInvalidCredentials) {
			status = bunker.StatusUnAuthorized
		}
		httpUtil.ErrorResponse(w, bunker.New(status).SetMessage(err.Error()))
		return
	}

	h.cookieManager.SetAccessTokenCookie(w, res.AccessToken)
	h.cookieManager.SetRefreshTokenCookie(w, res.RefreshToken)
	h.cookieManager.SetCSRFCookie(w, res.CSRFToken)

	httpUtil.SuccessResponse(w, res)
}

// RefreshToken menangani Refresh Token Rotation (Cookie-only)
func (h *authHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	tokenStr := ""
	if c, err := r.Cookie(cookie.RefreshTokenCookieName); err == nil && c.Value != "" {
		tokenStr = c.Value
	}

	if tokenStr == "" {
		httpUtil.ErrorResponse(w, bunker.New(bunker.StatusUnAuthorized).SetMessage(domain.ErrMissingRefreshToken.Error()))
		return
	}

	res, err := h.usecase.RefreshToken(r.Context(), tokenStr)
	if err != nil {
		status := bunker.StatusUnAuthorized
		if errors.Is(err, domain.ErrTokenReused) {
			status = bunker.StatusForbidden
			h.cookieManager.ClearAuthCookies(w)
		}
		httpUtil.ErrorResponse(w, bunker.New(status).SetMessage(err.Error()))
		return
	}

	h.cookieManager.SetAccessTokenCookie(w, res.AccessToken)
	h.cookieManager.SetRefreshTokenCookie(w, res.RefreshToken)
	h.cookieManager.SetCSRFCookie(w, res.CSRFToken)

	httpUtil.SuccessResponse(w, res)
}

// Logout mencabut sesi dan membersihkan cookies
func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	tokenStr := ""
	if c, err := r.Cookie(cookie.RefreshTokenCookieName); err == nil {
		tokenStr = c.Value
	}

	_ = h.usecase.Logout(r.Context(), tokenStr)
	h.cookieManager.ClearAuthCookies(w)

	httpUtil.SuccessResponse(w, map[string]string{
		"message": "berhasil keluar (logout)",
	})
}

// CSRFToken menghasilkan token CSRF awal untuk client (SPA / Web App)
func (h *authHandler) CSRFToken(w http.ResponseWriter, r *http.Request) {
	csrfToken, err := h.csrfManager.GenerateToken()
	if err != nil {
		httpUtil.ErrorResponse(w, bunker.New(bunker.StatusUnprocessableEntity).SetMessage(err.Error()))
		return
	}

	h.cookieManager.SetCSRFCookie(w, csrfToken)
	httpUtil.SuccessResponse(w, map[string]string{
		"csrf_token": csrfToken,
	})
}
