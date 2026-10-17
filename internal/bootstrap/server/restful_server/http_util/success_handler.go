package httpUtil

import (
	"encoding/json"
	"net/http"

	modelCommon "github.com/orchid-craft/orchid-auth-starter/internal/common/model"
	bunker "github.com/olympus-forge/bunker/errors"
)

func SuccessResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(modelCommon.Response{
		Success: true,
		Status:  http.StatusOK,
		Data:    data,
	})
}

// NewAppError mengembalikan salinan *bunker.ApplicationError yang independen dan aman
// digunakan secara konkuren. Ini mencegah race condition karena bunker.New mengembalikan
// pointer ke variabel global singleton (mis. bunker.Forbidden) yang termutasi jika SetMessage dipanggil.
func NewAppError(code bunker.AppErrorCode, message string) *bunker.ApplicationError {
	cp := *bunker.New(code)
	return cp.SetMessage(message)
}

func ErrorResponse(w http.ResponseWriter, err error) {
	bunkerErr, ok := err.(*bunker.ApplicationError)
	if !ok {
		bunkerErr = NewAppError(bunker.StatusUnprocessableEntity, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(bunkerErr.HttpStatusCode)
	_ = json.NewEncoder(w).Encode(modelCommon.Response{
		Success:      false,
		Status:       bunkerErr.HttpStatusCode,
		ErrorCode:    bunkerErr.ErrorCode,
		ErrorMessage: bunkerErr.Message,
	})
}
