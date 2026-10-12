package httpUtil

import (
	"net/http"

	bunker "github.com/yudhiana/bunker/errors"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	message := "route not found"
	err := bunker.New(bunker.StatusNotFound).SetMessage(message)
	ErrorResponse(w, err)
}
