package repository

import (
	"context"
	"sync"

	"github.com/orchid-craft/orchid-auth-starter/modules/user/domain"
	userModel "github.com/orchid-craft/orchid-auth-starter/modules/user/domain/models"
)

type memoryUserRepository struct {
	mu      sync.RWMutex
	users   map[string]*userModel.User
	byEmail map[string]string
}

func NewUserRepository() UserRepositoryInterface {
	return &memoryUserRepository{
		users:   make(map[string]*userModel.User),
		byEmail: make(map[string]string),
	}
}

func (r *memoryUserRepository) Create(ctx context.Context, user *userModel.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byEmail[user.Email]; exists {
		return domain.ErrUserAlreadyExists
	}

	r.users[user.ID] = user
	r.byEmail[user.Email] = user.ID
	return nil
}

func (r *memoryUserRepository) FindByEmail(ctx context.Context, email string) (*userModel.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exists := r.byEmail[email]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	user, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (r *memoryUserRepository) FindByID(ctx context.Context, id string) (*userModel.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (r *memoryUserRepository) Update(ctx context.Context, user *userModel.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return domain.ErrUserNotFound
	}

	r.users[user.ID] = user
	return nil
}

func (r *memoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return domain.ErrUserNotFound
	}

	delete(r.byEmail, user.Email)
	delete(r.users, id)
	return nil
}

func (r *memoryUserRepository) List(ctx context.Context) ([]*userModel.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*userModel.User, 0, len(r.users))
	for _, u := range r.users {
		list = append(list, u)
	}
	return list, nil
}
