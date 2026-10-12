package config

import "time"

type LocalConfig struct {
	AppName    string `env:"APP_NAME" envDefault:"web-api"`
	AppPort    string `env:"APP_PORT" envDefault:"8080"`
	AppHost    string `env:"APP_HOST" envDefault:"0.0.0.0"`
	AppVersion string `env:"APP_VERSION" envDefault:"1.0.0"`
	AppEnv     string `env:"APP_ENV" envDefault:"development"`
	DebugMode  bool   `env:"DEBUG_MODE" envDefault:"true"`

	// Security & Auth Configurations
	AccessTokenSecret      string `env:"ACCESS_TOKEN_SECRET" envDefault:"super-secret-access-token-key-change-in-production-min-32-chars!"`
	RefreshTokenSecret     string `env:"REFRESH_TOKEN_SECRET" envDefault:"super-secret-refresh-token-key-change-in-production-min-32-chars!"`
	CSRFSecret             string `env:"CSRF_SECRET" envDefault:"super-secret-csrf-key-change-in-production-min-32-chars!"`
	AccessTokenTTLMinutes  int    `env:"ACCESS_TOKEN_TTL_MINUTES" envDefault:"15"`
	RefreshTokenTTLDays    int    `env:"REFRESH_TOKEN_TTL_DAYS" envDefault:"7"`
	CookieSecure           bool   `env:"COOKIE_SECURE" envDefault:"false"`
	CookieDomain           string `env:"COOKIE_DOMAIN" envDefault:""`
}

func (c *LocalConfig) AccessTokenTTL() time.Duration {
	return time.Duration(c.AccessTokenTTLMinutes) * time.Minute
}

func (c *LocalConfig) RefreshTokenTTL() time.Duration {
	return time.Duration(c.RefreshTokenTTLDays) * 24 * time.Hour
}

func (c *LocalConfig) IsProduction() bool {
	return c.AppEnv == "production"
}
