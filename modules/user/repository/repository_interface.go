package repository

import (
	"context"

	userModel "github.com/orchid-craft/orchid-auth-starter/modules/user/domain/models"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *userModel.User) error
	FindByID(ctx context.Context, id string) (*userModel.User, error)
	FindByEmail(ctx context.Context, email string) (*userModel.User, error)
	Update(ctx context.Context, user *userModel.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*userModel.User, error)
}
