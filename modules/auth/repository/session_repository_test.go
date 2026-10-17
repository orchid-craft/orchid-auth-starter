package repository

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/orchid-craft/orchid-auth-starter/modules/auth/domain"
	authModel "github.com/orchid-craft/orchid-auth-starter/modules/auth/domain/models"
)

func newActiveSession(id, userID string) *authModel.RefreshTokenSession {
	return &authModel.RefreshTokenSession{
		ID:        id,
		UserID:    userID,
		TokenHash: "hash-" + id,
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
	}
}

// Rotasi refresh token tidak boleh meloloskan dua request paralel dengan token yang sama.
// Jalankan dengan: go test -race ./modules/auth/repository/...
func TestConsume_OnlyOneConcurrentCallerWins(t *testing.T) {
	repo := NewSessionRepository()
	ctx := context.Background()

	if err := repo.Save(ctx, newActiveSession("jti-1", "user-1")); err != nil {
		t.Fatalf("Save gagal: %v", err)
	}

	const workers = 50
	var success, reused atomic.Int32
	var wg sync.WaitGroup
	start := make(chan struct{})

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, err := repo.Consume(ctx, "jti-1")
			switch {
			case err == nil:
				success.Add(1)
			case errors.Is(err, domain.ErrTokenReused):
				reused.Add(1)
			default:
				t.Errorf("error tak terduga dari Consume: %v", err)
			}
		}()
	}

	close(start)
	wg.Wait()

	if got := success.Load(); got != 1 {
		t.Fatalf("diharapkan tepat 1 Consume sukses, didapat %d", got)
	}
	if got := reused.Load(); got != workers-1 {
		t.Fatalf("diharapkan %d Consume ditolak sebagai reuse, didapat %d", workers-1, got)
	}
}

func TestConsume_SecondCallIsReuse(t *testing.T) {
	repo := NewSessionRepository()
	ctx := context.Background()
	_ = repo.Save(ctx, newActiveSession("jti-1", "user-1"))

	if _, err := repo.Consume(ctx, "jti-1"); err != nil {
		t.Fatalf("Consume pertama harus sukses, didapat: %v", err)
	}
	if _, err := repo.Consume(ctx, "jti-1"); !errors.Is(err, domain.ErrTokenReused) {
		t.Fatalf("Consume kedua harus ErrTokenReused, didapat: %v", err)
	}
}

func TestConsume_UnknownSession(t *testing.T) {
	repo := NewSessionRepository()

	if _, err := repo.Consume(context.Background(), "tidak-ada"); !errors.Is(err, domain.ErrTokenInvalid) {
		t.Fatalf("sesi tidak dikenal harus ErrTokenInvalid, didapat: %v", err)
	}
}

func TestConsume_RevokedByLogoutIsReuse(t *testing.T) {
	repo := NewSessionRepository()
	ctx := context.Background()
	_ = repo.Save(ctx, newActiveSession("jti-1", "user-1"))
	_ = repo.Revoke(ctx, "jti-1")

	if _, err := repo.Consume(ctx, "jti-1"); !errors.Is(err, domain.ErrTokenReused) {
		t.Fatalf("sesi yang sudah dicabut harus ErrTokenReused, didapat: %v", err)
	}
}

// FindByID harus mengembalikan salinan: memutasi hasilnya tidak boleh mengubah state repository.
func TestFindByID_ReturnsCopy(t *testing.T) {
	repo := NewSessionRepository()
	ctx := context.Background()
	_ = repo.Save(ctx, newActiveSession("jti-1", "user-1"))

	found, err := repo.FindByID(ctx, "jti-1")
	if err != nil {
		t.Fatalf("FindByID gagal: %v", err)
	}
	found.Revoked = true
	found.UserID = "dimodifikasi"

	again, _ := repo.FindByID(ctx, "jti-1")
	if again.Revoked || again.UserID != "user-1" {
		t.Fatalf("FindByID mengekspos state internal: %+v", again)
	}
}

// Save harus menyimpan salinan: mutasi pointer milik pemanggil setelah Save tidak boleh bocor ke repository.
func TestSave_StoresCopy(t *testing.T) {
	repo := NewSessionRepository()
	ctx := context.Background()

	session := newActiveSession("jti-1", "user-1")
	_ = repo.Save(ctx, session)
	session.Revoked = true

	found, _ := repo.FindByID(ctx, "jti-1")
	if found.Revoked {
		t.Fatal("Save menyimpan pointer pemanggil, bukan salinan")
	}
}

func TestRevokeAllForUser_OnlyAffectsThatUser(t *testing.T) {
	repo := NewSessionRepository()
	ctx := context.Background()
	_ = repo.Save(ctx, newActiveSession("a1", "user-1"))
	_ = repo.Save(ctx, newActiveSession("a2", "user-1"))
	_ = repo.Save(ctx, newActiveSession("b1", "user-2"))

	_ = repo.RevokeAllForUser(ctx, "user-1")

	for _, id := range []string{"a1", "a2"} {
		if s, _ := repo.FindByID(ctx, id); !s.Revoked {
			t.Errorf("sesi %s milik user-1 seharusnya revoked", id)
		}
	}
	if s, _ := repo.FindByID(ctx, "b1"); s.Revoked {
		t.Error("sesi milik user-2 tidak boleh ikut revoked")
	}
}
