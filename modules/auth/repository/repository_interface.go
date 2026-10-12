package repository

import (
	"context"

	authModel "github.com/yudhiana/web-api/modules/auth/domain/models"
)

type SessionRepositoryInterface interface {
	Save(ctx context.Context, session *authModel.RefreshTokenSession) error
	FindByID(ctx context.Context, id string) (*authModel.RefreshTokenSession, error)
	Revoke(ctx context.Context, id string) error
	RevokeAllForUser(ctx context.Context, userID string) error
}
