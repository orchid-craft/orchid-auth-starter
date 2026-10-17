package container

import (
	"context"
	"errors"
	"testing"

	"github.com/orchid-craft/orchid-auth-starter/config"
	userDomain "github.com/orchid-craft/orchid-auth-starter/modules/user/domain"
)

func TestSeedDefaultUsers_OnlyInDevelopmentAndTest(t *testing.T) {
	cases := []struct {
		env        string
		wantSeeded bool
	}{
		{"development", true},
		{"test", true},
		{"production", false},
		{"prod", false},
		{"staging", false},
		{"", false},
	}

	for _, tc := range cases {
		t.Run("APP_ENV="+tc.env, func(t *testing.T) {
			di, err := NewApplicationContainer(&config.LocalConfig{AppEnv: tc.env})
			if err != nil {
				t.Fatalf("NewApplicationContainer gagal: %v", err)
			}

			_, err = di.UserUsecase.GetByEmail(context.Background(), "alice@example.com")
			seeded := err == nil

			if !seeded && !errors.Is(err, userDomain.ErrUserNotFound) {
				t.Fatalf("error tak terduga: %v", err)
			}
			if seeded != tc.wantSeeded {
				t.Fatalf("APP_ENV=%q: akun admin demo seeded=%v, diharapkan %v", tc.env, seeded, tc.wantSeeded)
			}
		})
	}
}
