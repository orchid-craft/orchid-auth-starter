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

func (r *memorySessionRepository) Save(ctx context.Context, session *authModel.RefreshTokenSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[session.ID] = session
	return nil
}

func (r *memorySessionRepository) FindByID(ctx context.Context, id string) (*authModel.RefreshTokenSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, exists := r.sessions[id]
	if !exists {
		return nil, domain.ErrTokenInvalid
	}
	return session, nil
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
