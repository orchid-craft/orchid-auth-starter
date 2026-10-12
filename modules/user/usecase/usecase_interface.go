package usecase

import (
	"context"

	reqModel "github.com/yudhiana/web-api/modules/user/delivery/models/request"
	resModel "github.com/yudhiana/web-api/modules/user/delivery/models/response"
	userModel "github.com/yudhiana/web-api/modules/user/domain/models"
)

type UserUsecaseInterface interface {
	Register(ctx context.Context, req reqModel.CreateUserRequest) (*resModel.UserResponse, error)
	GetByID(ctx context.Context, id string) (*resModel.UserResponse, error)
	GetByEmail(ctx context.Context, email string) (*userModel.User, error)
	Update(ctx context.Context, id string, req reqModel.UpdateUserRequest) (*resModel.UserResponse, error)
	List(ctx context.Context) ([]*resModel.UserResponse, error)
}
