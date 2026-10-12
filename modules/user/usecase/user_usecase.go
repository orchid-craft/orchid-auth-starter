package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	reqModel "github.com/yudhiana/web-api/modules/user/delivery/models/request"
	resModel "github.com/yudhiana/web-api/modules/user/delivery/models/response"
	userModel "github.com/yudhiana/web-api/modules/user/domain/models"
	"github.com/yudhiana/web-api/modules/user/repository"
)

type userUsecase struct {
	repo repository.UserRepositoryInterface
}

func NewUserUsecase(repo repository.UserRepositoryInterface) UserUsecaseInterface {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) Register(ctx context.Context, req reqModel.CreateUserRequest) (*resModel.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("gagal mengenkripsi kata sandi: %w", err)
	}

	role := req.Role
	if role != "admin" {
		role = "user"
	}

	now := time.Now()
	user := &userModel.User{
		ID:           uuid.NewString(),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         role,
		DisplayName:  req.DisplayName,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

func (u *userUsecase) GetByID(ctx context.Context, id string) (*resModel.UserResponse, error) {
	user, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toUserResponse(user), nil
}

func (u *userUsecase) GetByEmail(ctx context.Context, email string) (*userModel.User, error) {
	return u.repo.FindByEmail(ctx, email)
}

func (u *userUsecase) Update(ctx context.Context, id string, req reqModel.UpdateUserRequest) (*resModel.UserResponse, error) {
	user, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.DisplayName = req.DisplayName
	user.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

func (u *userUsecase) List(ctx context.Context) ([]*resModel.UserResponse, error) {
	users, err := u.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*resModel.UserResponse, len(users))
	for i, user := range users {
		results[i] = toUserResponse(user)
	}
	return results, nil
}

func toUserResponse(u *userModel.User) *resModel.UserResponse {
	return &resModel.UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		Role:        u.Role,
		DisplayName: u.DisplayName,
		CreatedAt:   u.CreatedAt,
	}
}
