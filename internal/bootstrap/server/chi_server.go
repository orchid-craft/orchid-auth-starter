package server

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/olympus-forge/logos"
	"github.com/yudhiana/web-api/constants"
	httpUtil "github.com/yudhiana/web-api/internal/bootstrap/server/restful_server/http_util"
	modelCommon "github.com/yudhiana/web-api/internal/common/model"
)

type chiApplication struct {
	app *chi.Mux
}

var chiAppInstance *chiApplication
var onceApp sync.Once

// NewChiApplication creates a fresh chi.Mux router with standard middlewares (P2-1)
func NewChiApplication() *chi.Mux {
	log := logos.NewLogger()
	log.Info("Initialize chi application instance...")

	app := chi.NewRouter()
	app.Use(middleware.RequestID)

	// Propagasi Request ID ke context untuk logging dan usecase, serta sematkan ke response header (P2-4)
	app.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			if reqID == "" {
				reqID = r.Header.Get(constants.HeaderRequestID)
			}
			if reqID != "" {
				w.Header().Set(constants.HeaderRequestID, reqID)
			}
			ctx := modelCommon.SetRequestContext(r.Context(), &modelCommon.RequestContext{
				RequestID: reqID,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)
	app.NotFound(httpUtil.NotFoundHandler)

	return app
}

func GetChiApplication() *chi.Mux {
	onceApp.Do(func() {
		chiAppInstance = &chiApplication{
			app: NewChiApplication(),
		}
	})

	return chiAppInstance.app
}
