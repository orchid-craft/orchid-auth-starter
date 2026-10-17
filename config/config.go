package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// GetLocalConfig loads configuration from environment variables and .env file
func GetLocalConfig() *LocalConfig {
	// Attempt loading .env from current directory or parent directories
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../../.env")

	var localConfig LocalConfig
	if err := env.Parse(&localConfig); err != nil {
		log.Fatalf("failed to parse environment configurations: %v", err)
	}

	// P1-4 / P0: validasi fail-closed. Hanya APP_ENV "development"/"test" yang boleh memakai secret default;
	// selain itu server menolak start jika secret lemah/placeholder/kembar atau COOKIE_SECURE=false.
	if err := localConfig.Validate(); err != nil {
		log.Fatalf("keamanan kritis: %v", err)
	}

	return &localConfig
}
