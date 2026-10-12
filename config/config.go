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

	return &localConfig
}
