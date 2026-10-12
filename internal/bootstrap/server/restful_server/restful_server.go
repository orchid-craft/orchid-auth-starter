package restfulServer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/olympus-forge/logos"
	"github.com/yudhiana/web-api/config"
	"github.com/yudhiana/web-api/internal/bootstrap/container"
)

type Server struct {
	cfg *config.LocalConfig
	app *chi.Mux
	log *logos.LogEntry
}

func NewServer(c *container.Container) *Server {
	srv := &Server{
		cfg: c.GetConfig(),
		app: c.GetApp(),
		log: c.Log,
	}

	srv.setupRoutes(c)
	return srv
}

func (s *Server) setupRoutes(c *container.Container) {
	routeSetups := []RouteSetup{
		SetupCommonRoutes,
		SetupAuthRoutes,
		SetupUserRoutes,
	}

	SetupAllRoutes(s.app, c, routeSetups...)
}

func (s *Server) PrintRouteBanner() {
	routesList := make([]struct{ Method, Path string }, 0)

	_ = chi.Walk(s.app, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		// Clean up route formatting
		cleanRoute := strings.ReplaceAll(route, "/*/", "/")
		routesList = append(routesList, struct{ Method, Path string }{
			Method: method,
			Path:   cleanRoute,
		})
		return nil
	})

	fmt.Println("========================================================================")
	fmt.Printf("🚀  ORCHID STARTER WEB API (go-chi/chi/v5)\n")
	fmt.Printf("📡  Port Server : :%s\n", s.cfg.AppPort)
	fmt.Printf("🌐  Host Server : %s\n", s.cfg.AppHost)
	fmt.Printf("🔒  Environment : %s\n", s.cfg.AppEnv)
	fmt.Printf("🛠️   Debug Mode  : %t\n", s.cfg.DebugMode)
	fmt.Printf("📋  Total Routes: %d Rute Terdaftar\n", len(routesList))
	fmt.Println("------------------------------------------------------------------------")
	fmt.Println("  DAFTAR ROUTER ENDPOINT:")
	for _, route := range routesList {
		color := "\033[0;32m" // default green (GET)
		switch route.Method {
		case "POST":
			color = "\033[0;34m" // blue
		case "PUT":
			color = "\033[1;33m" // yellow
		case "DELETE":
			color = "\033[0;31m" // red
		}
		reset := "\033[0m"
		fmt.Printf("  %s[%-6s]%s %s\n", color, route.Method, reset, route.Path)
	}
	fmt.Println("========================================================================")
	if s.cfg.DebugMode {
		fmt.Printf("💡 Debug Logger AKTIF: RequestID dan log HTTP request/response dicatat secara real-time.\n\n")
	}
}

func (s *Server) Run() error {
	cfg := s.cfg

	if cfg.DebugMode {
		s.PrintRouteBanner()
	}

	s.log.Info(
		"Server starting",
		"version", cfg.AppVersion,
		"host", cfg.AppHost,
		"port", cfg.AppPort,
	)

	address := fmt.Sprintf("%s:%s", cfg.AppHost, cfg.AppPort)
	srv := &http.Server{
		Addr:    address,
		Handler: s.app,
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		s.log.Info("🔴Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			s.log.Error("🔴Forced shutdown", "err", err)
		}

		s.log.Info("🔴Server exited properly")
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	return nil
}
