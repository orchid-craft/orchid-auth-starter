package config

import (
	"reflect"
	"strings"
	"testing"
)

// validProductionConfig mengembalikan konfigurasi production yang lolos Validate().
func validProductionConfig() LocalConfig {
	return LocalConfig{
		AppEnv:             "production",
		AccessTokenSecret:  strings.Repeat("a", 40),
		RefreshTokenSecret: strings.Repeat("b", 40),
		CSRFSecret:         strings.Repeat("c", 40),
		CookieSecure:       true,
	}
}

func defaultsConfig(env string) LocalConfig {
	return LocalConfig{
		AppEnv:             env,
		AccessTokenSecret:  defaultAccessTokenSecret,
		RefreshTokenSecret: defaultRefreshTokenSecret,
		CSRFSecret:         defaultCSRFSecret,
		CookieSecure:       false,
	}
}

func TestValidate_DevelopmentAndTestAllowDefaults(t *testing.T) {
	for _, env := range []string{"development", "test", "Development", " TEST "} {
		cfg := defaultsConfig(env)
		if err := cfg.Validate(); err != nil {
			t.Errorf("APP_ENV=%q seharusnya boleh memakai default, tetapi error: %v", env, err)
		}
		if !cfg.IsDevelopment() || cfg.IsProduction() {
			t.Errorf("APP_ENV=%q seharusnya dikenali sebagai development", env)
		}
	}
}

// Inti perbaikan fail-open: env kosong / salah ketik TIDAK boleh membuat secret default diterima.
func TestValidate_FailsClosedForAnythingNotDevelopment(t *testing.T) {
	for _, env := range []string{"", "prod", "Production", "production", "staging", "dev", "local"} {
		cfg := defaultsConfig(env)
		err := cfg.Validate()
		if err == nil {
			t.Errorf("APP_ENV=%q dengan secret default seharusnya DITOLAK", env)
			continue
		}
		if !cfg.IsProduction() {
			t.Errorf("APP_ENV=%q seharusnya diperlakukan sebagai production", env)
		}
		for _, name := range []string{"ACCESS_TOKEN_SECRET", "REFRESH_TOKEN_SECRET", "CSRF_SECRET", "COOKIE_SECURE"} {
			if !strings.Contains(err.Error(), name) {
				t.Errorf("APP_ENV=%q: pesan error seharusnya menyebut %s, didapat: %v", env, name, err)
			}
		}
	}
}

func TestValidate_ValidProductionConfigPasses(t *testing.T) {
	for _, env := range []string{"production", "prod", "staging", "Production"} {
		cfg := validProductionConfig()
		cfg.AppEnv = env
		if err := cfg.Validate(); err != nil {
			t.Errorf("APP_ENV=%q dengan konfigurasi kuat seharusnya lolos, tetapi error: %v", env, err)
		}
	}
}

// Regresi: docker-compose lama memakai default + suffix "!prod" yang lolos perbandingan string persis.
func TestValidate_RejectsDefaultWithSuffix(t *testing.T) {
	cfg := validProductionConfig()
	cfg.AccessTokenSecret = defaultAccessTokenSecret + "prod"
	cfg.RefreshTokenSecret = defaultRefreshTokenSecret + "prod"
	cfg.CSRFSecret = defaultCSRFSecret + "prod"

	err := cfg.Validate()
	if err == nil {
		t.Fatal("secret default + suffix seharusnya ditolak")
	}
	if strings.Count(err.Error(), "default/placeholder") != 3 {
		t.Errorf("ketiga secret seharusnya ditandai placeholder, didapat: %v", err)
	}
}

func TestValidate_RejectsShortSecret(t *testing.T) {
	cfg := validProductionConfig()
	cfg.CSRFSecret = "terlalu-pendek"

	err := cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "CSRF_SECRET wajib diatur minimal 32 karakter") {
		t.Fatalf("secret pendek seharusnya ditolak, didapat: %v", err)
	}
}

func TestValidate_RejectsDuplicateSecrets(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*LocalConfig)
		want   string
	}{
		{"access==refresh", func(c *LocalConfig) { c.RefreshTokenSecret = c.AccessTokenSecret }, "ACCESS_TOKEN_SECRET dan REFRESH_TOKEN_SECRET"},
		{"access==csrf", func(c *LocalConfig) { c.CSRFSecret = c.AccessTokenSecret }, "ACCESS_TOKEN_SECRET dan CSRF_SECRET"},
		{"refresh==csrf", func(c *LocalConfig) { c.CSRFSecret = c.RefreshTokenSecret }, "REFRESH_TOKEN_SECRET dan CSRF_SECRET"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := validProductionConfig()
			tc.mutate(&cfg)

			err := cfg.Validate()
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("secret kembar (%s) seharusnya ditolak, didapat: %v", tc.name, err)
			}
		})
	}
}

func TestValidate_RequiresCookieSecureOutsideDevelopment(t *testing.T) {
	cfg := validProductionConfig()
	cfg.CookieSecure = false

	err := cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "COOKIE_SECURE") {
		t.Fatalf("COOKIE_SECURE=false di production seharusnya ditolak, didapat: %v", err)
	}
}

// Menjaga konstanta di validate.go tetap sinkron dengan envDefault pada struct.
func TestDefaultSecretsMatchStructTags(t *testing.T) {
	typ := reflect.TypeOf(LocalConfig{})
	cases := map[string]string{
		"AccessTokenSecret":  defaultAccessTokenSecret,
		"RefreshTokenSecret": defaultRefreshTokenSecret,
		"CSRFSecret":         defaultCSRFSecret,
	}

	for field, want := range cases {
		f, ok := typ.FieldByName(field)
		if !ok {
			t.Fatalf("field %s tidak ditemukan", field)
		}
		if got := f.Tag.Get("envDefault"); got != want {
			t.Errorf("envDefault %s (%q) tidak sama dengan konstanta di validate.go (%q)", field, got, want)
		}
	}
}
