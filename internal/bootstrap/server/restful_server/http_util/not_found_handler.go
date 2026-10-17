package httpUtil

import (
	"net/http"

	bunker "github.com/olympus-forge/bunker/errors"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, NewAppError(bunker.StatusNotFound, "route not found"))
}
