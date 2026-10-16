package httpUtil

import (
	"encoding/json"
	"net/http"

	modelCommon "github.com/orchid-craft/orchid-auth-starter/internal/common/model"
	bunker "github.com/yudhiana/bunker/errors"
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

func ErrorResponse(w http.ResponseWriter, err error) {
	bunkerErr, ok := err.(*bunker.ApplicationError)
	if !ok {
		bunkerErr = bunker.New(bunker.StatusUnprocessableEntity).SetMessage(err.Error())
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
