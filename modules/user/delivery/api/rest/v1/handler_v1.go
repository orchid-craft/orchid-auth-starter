package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	bunker "github.com/olympus-forge/bunker/errors"
	httpUtil "github.com/orchid-craft/orchid-auth-starter/internal/bootstrap/server/restful_server/http_util"
	reqModel "github.com/orchid-craft/orchid-auth-starter/modules/user/delivery/models/request"
	"github.com/orchid-craft/orchid-auth-starter/modules/user/domain"
	"github.com/orchid-craft/orchid-auth-starter/modules/user/usecase"
	"github.com/orchid-craft/orchid-auth-starter/pkg/validator"
	"github.com/orchid-craft/orchid-auth-starter/security"
)

type userHandler struct {
	usecase usecase.UserUsecaseInterface
}

func NewUserHandler(u usecase.UserUsecaseInterface) *userHandler {
	return &userHandler{usecase: u}
}

// GetProfile mengembalikan profil user yang sedang login dari context claims
func (h *userHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.GetClaimsContext(r.Context())
	if !ok || claims == nil {
		err := httpUtil.NewAppError(bunker.StatusUnAuthorized, "claims pengguna tidak ditemukan di context")
		httpUtil.ErrorResponse(w, err)
		return
	}

	user, err := h.usecase.GetByID(r.Context(), claims.UserID)
	if err != nil {
		status := bunker.StatusUnprocessableEntity
		if errors.Is(err, domain.ErrUserNotFound) {
			status = bunker.StatusNotFound
		}
		httpUtil.ErrorResponse(w, httpUtil.NewAppError(status, err.Error()))
		return
	}

	httpUtil.SuccessResponse(w, user)
}

// UpdateProfile memperbarui profil user yang sedang login
func (h *userHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.GetClaimsContext(r.Context())
	if !ok || claims == nil {
		err := httpUtil.NewAppError(bunker.StatusUnAuthorized, "claims pengguna tidak ditemukan di context")
		httpUtil.ErrorResponse(w, err)
		return
	}

	var req reqModel.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errBunker := httpUtil.NewAppError(bunker.StatusBadRequest, "format request payload tidak valid")
		httpUtil.ErrorResponse(w, errBunker)
		return
	}

	if err := validator.ValidateStruct(&req); err != nil {
		httpUtil.ErrorResponse(w, httpUtil.NewAppError(bunker.StatusBadRequest, err.Error()))
		return
	}

	updated, err := h.usecase.Update(r.Context(), claims.UserID, req)
	if err != nil {
		status := bunker.StatusUnprocessableEntity
		if errors.Is(err, domain.ErrUserNotFound) {
			status = bunker.StatusNotFound
		}
		httpUtil.ErrorResponse(w, httpUtil.NewAppError(status, err.Error()))
		return
	}

	httpUtil.SuccessResponse(w, updated)
}

// ListUsers mengembalikan daftar seluruh user (hanya untuk admin / authorized)
func (h *userHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.usecase.List(r.Context())
	if err != nil {
		httpUtil.ErrorResponse(w, httpUtil.NewAppError(bunker.StatusUnprocessableEntity, err.Error()))
		return
	}

	httpUtil.SuccessResponse(w, users)
}
