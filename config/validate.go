package config

import (
	"errors"
	"fmt"
	"strings"
)

// minSecretLength adalah panjang minimum (karakter) untuk setiap secret di luar development/test.
const minSecretLength = 32

// Nilai default secret untuk development/test. HARUS identik dengan envDefault pada LocalConfig
// (dijaga oleh TestDefaultSecretsMatchStructTags).
const (
	defaultAccessTokenSecret  = "super-secret-access-token-key-change-in-production-min-32-chars!"
	defaultRefreshTokenSecret = "super-secret-refresh-token-key-change-in-production-min-32-chars!"
	defaultCSRFSecret         = "super-secret-csrf-key-change-in-production-min-32-chars!"
)

// Validate memeriksa keamanan konfigurasi dengan prinsip fail-closed.
//
// Development/test (APP_ENV eksplisit "development" atau "test") boleh memakai default.
// Lingkungan lain, termasuk APP_ENV kosong atau salah ketik, WAJIB memenuhi:
//   - ketiga secret minimal 32 karakter dan bukan default/placeholder
//   - ACCESS_TOKEN_SECRET, REFRESH_TOKEN_SECRET, dan CSRF_SECRET saling berbeda (isolasi dual-secret)
//   - COOKIE_SECURE=true
//
// Semua pelanggaran dikumpulkan dalam satu error agar operator bisa memperbaikinya sekaligus.
func (c *LocalConfig) Validate() error {
	if c.IsDevelopment() {
		return nil
	}

	var errs []error

	secrets := []struct{ name, value string }{
		{"ACCESS_TOKEN_SECRET", c.AccessTokenSecret},
		{"REFRESH_TOKEN_SECRET", c.RefreshTokenSecret},
		{"CSRF_SECRET", c.CSRFSecret},
	}

	for _, s := range secrets {
		switch {
		case len(s.value) < minSecretLength:
			errs = append(errs, fmt.Errorf("%s wajib diatur minimal %d karakter", s.name, minSecretLength))
		case isPlaceholderSecret(s.value):
			errs = append(errs, fmt.Errorf("%s masih berupa nilai default/placeholder; ganti dengan secret acak yang unik", s.name))
		}
	}

	for i := range secrets {
		for j := i + 1; j < len(secrets); j++ {
			if secrets[i].value != "" && secrets[i].value == secrets[j].value {
				errs = append(errs, fmt.Errorf("%s dan %s tidak boleh sama (isolasi dual-secret)", secrets[i].name, secrets[j].name))
			}
		}
	}

	if !c.CookieSecure {
		errs = append(errs, errors.New("COOKIE_SECURE wajib true di luar development/test"))
	}

	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("konfigurasi tidak aman untuk APP_ENV=%q (hanya \"development\" dan \"test\" yang boleh memakai default): %w",
		c.AppEnv, errors.Join(errs...))
}

// isPlaceholderSecret mendeteksi secret default maupun varian turunannya (mis. default + suffix "!prod"),
// yang lolos jika hanya dibandingkan dengan string default secara persis.
func isPlaceholderSecret(secret string) bool {
	switch secret {
	case defaultAccessTokenSecret, defaultRefreshTokenSecret, defaultCSRFSecret:
		return true
	}

	lower := strings.ToLower(secret)
	return strings.Contains(lower, "change-in-production") || strings.Contains(lower, "super-secret")
}
