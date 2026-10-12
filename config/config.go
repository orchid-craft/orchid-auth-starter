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

	// P1-4: Validasi secret di lingkungan production demi keamanan
	if localConfig.IsProduction() {
		const defaultAccess = "super-secret-access-token-key-change-in-production-min-32-chars!"
		const defaultRefresh = "super-secret-refresh-token-key-change-in-production-min-32-chars!"
		const defaultCSRF = "super-secret-csrf-key-change-in-production-min-32-chars!"

		if localConfig.AccessTokenSecret == defaultAccess || len(localConfig.AccessTokenSecret) < 32 {
			log.Fatalf("keamanan kritis: ACCESS_TOKEN_SECRET wajib diatur minimal 32 karakter dan tidak boleh menggunakan nilai default pada production")
		}
		if localConfig.RefreshTokenSecret == defaultRefresh || len(localConfig.RefreshTokenSecret) < 32 {
			log.Fatalf("keamanan kritis: REFRESH_TOKEN_SECRET wajib diatur minimal 32 karakter dan tidak boleh menggunakan nilai default pada production")
		}
		if localConfig.CSRFSecret == defaultCSRF || len(localConfig.CSRFSecret) < 32 {
			log.Fatalf("keamanan kritis: CSRF_SECRET wajib diatur minimal 32 karakter dan tidak boleh menggunakan nilai default pada production")
		}
	}

	return &localConfig
}
