package container

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/olympus-forge/logos"
	"github.com/yudhiana/web-api/config"
	"github.com/yudhiana/web-api/internal/bootstrap/server"
)

type Container struct {
	App *chi.Mux
	Cfg *config.LocalConfig
	DI  *DirectInjection
	Log *logos.LogEntry
}

// NewContainer creates a new application container with configuration, chi router, and DI
func NewContainer() (*Container, error) {
	logger := logos.NewLogger()
	logger.Info("Initializing application container...")

	// 1. Load configuration
	cfg := config.GetLocalConfig()
	if cfg == nil {
		return nil, fmt.Errorf("failed to load configuration")
	}
	logger.Info("Configuration loaded successfully")

	// 2. Initialize dependency injection
	di, err := NewApplicationContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize dependencies: %w", err)
	}
	logger.Info("Dependencies initialized successfully")

	// 3. Initialize fresh Chi router (P2-1)
	app := server.NewChiApplication()
	if app == nil {
		return nil, fmt.Errorf("failed to initialize chi application")
	}
	logger.Info("Chi application router initialized successfully")

	container := &Container{
		App: app,
		Cfg: cfg,
		DI:  di,
		Log: logger,
	}

	logger.Info("Application container initialized successfully")
	return container, nil
}

// Close gracefully shuts down the container and its resources
func (c *Container) Close() error {
	logger := logos.NewLogger()
	logger.Info("🔴Shutting down application container...")

	if err := c.DI.Close(); err != nil {
		return fmt.Errorf("container shutdown completed with errors: %w", err)
	}

	logger.Info("🔴Application container shut down successfully")
	return nil
}

func (c *Container) GetConfig() *config.LocalConfig {
	return c.Cfg
}

func (c *Container) GetApp() *chi.Mux {
	return c.App
}

func (c *Container) GetDI() *DirectInjection {
	return c.DI
}
