package server

import (
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/olympus-forge/logos"
	httpUtil "github.com/yudhiana/web-api/internal/bootstrap/server/restful_server/http_util"
)

type chiApplication struct {
	app *chi.Mux
}

var chiAppInstance *chiApplication
var onceApp sync.Once

func GetChiApplication() *chi.Mux {
	onceApp.Do(func() {
		log := logos.NewLogger()
		log.Info("Initialize chi application instance...")

		app := chi.NewRouter()
		app.Use(middleware.RequestID)
		app.Use(middleware.Logger)
		app.Use(middleware.Recoverer)
		app.NotFound(httpUtil.NotFoundHandler)

		chiAppInstance = &chiApplication{
			app: app,
		}
	})

	return chiAppInstance.app
}
