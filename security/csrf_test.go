package security

import (
	"testing"
)

func TestCSRFManager(t *testing.T) {
	mgr := NewCSRFManager("test-csrf-secret-minimum-32-chars-long")

	token, err := mgr.GenerateToken()
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// 1. Valid matching token
	if err := mgr.ValidateToken(token, token); err != nil {
		t.Errorf("Valid token rejected: %v", err)
	}

	// 2. Missing token
	if err := mgr.ValidateToken("", token); err == nil {
		t.Error("Harus error saat header token kosong")
	}

	// 3. Mismatched token
	if err := mgr.ValidateToken("some-other-token", token); err == nil {
		t.Error("Harus error saat token tidak sama")
	}

	// 4. Tampered signature
	fakeToken := "validnonce.fakesignature"
	if err := mgr.ValidateToken(fakeToken, fakeToken); err == nil {
		t.Error("Harus error saat signature dipalsukan")
	}
}
