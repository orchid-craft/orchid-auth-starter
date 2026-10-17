package httpUtil

import (
	"net/http"

	bunker "github.com/olympus-forge/bunker/errors"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	message := "route not found"
	err := bunker.New(bunker.StatusNotFound).SetMessage(message)
	ErrorResponse(w, err)
}
