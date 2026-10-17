package repository

import (
	"context"
	"sync"

	"github.com/orchid-craft/orchid-auth-starter/modules/auth/domain"
	authModel "github.com/orchid-craft/orchid-auth-starter/modules/auth/domain/models"
)

type memorySessionRepository struct {
	mu       sync.RWMutex
	sessions map[string]*authModel.RefreshTokenSession // key: session.ID (token JTI)
}

func NewSessionRepository() SessionRepositoryInterface {
	return &memorySessionRepository{
		sessions: make(map[string]*authModel.RefreshTokenSession),
	}
}

// Save menyimpan SALINAN sesi, sehingga mutasi pada pointer milik pemanggil
// tidak mengubah state repository di luar lock.
func (r *memorySessionRepository) Save(ctx context.Context, session *authModel.RefreshTokenSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	stored := *session
	r.sessions[session.ID] = &stored
	return nil
}

// FindByID mengembalikan SALINAN sesi (bukan pointer ke data internal) agar aman dari data race.
func (r *memorySessionRepository) FindByID(ctx context.Context, id string) (*authModel.RefreshTokenSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, exists := r.sessions[id]
	if !exists {
		return nil, domain.ErrTokenInvalid
	}

	snapshot := *session
	return &snapshot, nil
}

// Consume memeriksa dan mencabut sesi dalam SATU critical section.
// Ini menutup celah check-then-act pada Refresh Token Rotation: dua request paralel dengan
// refresh token yang sama tidak bisa sama-sama lolos; yang kedua mendapat ErrTokenReused.
func (r *memorySessionRepository) Consume(ctx context.Context, id string) (*authModel.RefreshTokenSession, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	session, exists := r.sessions[id]
	if !exists {
		return nil, domain.ErrTokenInvalid
	}

	snapshot := *session
	if session.Revoked {
		return &snapshot, domain.ErrTokenReused
	}

	session.Revoked = true
	snapshot.Revoked = true
	return &snapshot, nil
}

func (r *memorySessionRepository) Revoke(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	session, exists := r.sessions[id]
	if !exists {
		return domain.ErrTokenInvalid
	}

	session.Revoked = true
	return nil
}

func (r *memorySessionRepository) RevokeAllForUser(ctx context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, session := range r.sessions {
		if session.UserID == userID {
			session.Revoked = true
		}
	}
	return nil
}
