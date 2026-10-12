package security

import (
	"testing"
	"time"
)

func TestJWTManager_GenerateAndValidate(t *testing.T) {
	cfg := JWTConfig{
		AccessTokenSecret:  "test-access-token-secret-minimum-32-chars-long",
		RefreshTokenSecret: "test-refresh-token-secret-minimum-32-chars-long",
		AccessTokenTTL:     15 * time.Minute,
		RefreshTokenTTL:    7 * 24 * time.Hour,
	}
	mgr := NewJWTManager(cfg)

	// 1. Access Token
	accToken, err := mgr.GenerateAccessToken("user-123", "user@example.com", "admin")
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	claims, err := mgr.ValidateAccessToken(accToken)
	if err != nil {
		t.Fatalf("ValidateAccessToken failed: %v", err)
	}
	if claims.UserID != "user-123" || claims.Email != "user@example.com" || claims.Role != "admin" {
		t.Errorf("Claims mismatch: %+v", claims)
	}

	// 2. Secret Separation: Verifikasi refresh secret tidak bisa memvalidasi access token
	otherMgr := NewJWTManager(JWTConfig{
		AccessTokenSecret:  "wrong-secret-key-32-chars-long-minimum!",
		RefreshTokenSecret: cfg.RefreshTokenSecret,
		AccessTokenTTL:     cfg.AccessTokenTTL,
		RefreshTokenTTL:    cfg.RefreshTokenTTL,
	})
	if _, err := otherMgr.ValidateAccessToken(accToken); err == nil {
		t.Error("Validasi harus gagal saat secret access token salah")
	}

	// 3. Refresh Token
	refToken, tokenID, expiresAt, err := mgr.GenerateRefreshToken("user-123")
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}
	if tokenID == "" {
		t.Error("Token ID (JTI) tidak boleh kosong")
	}
	if expiresAt.Before(time.Now()) {
		t.Error("ExpiresAt tidak boleh lampau")
	}

	refClaims, err := mgr.ValidateRefreshToken(refToken)
	if err != nil {
		t.Fatalf("ValidateRefreshToken failed: %v", err)
	}
	if refClaims.UserID != "user-123" || refClaims.TokenID != tokenID {
		t.Errorf("Refresh claims mismatch: %+v", refClaims)
	}
}
