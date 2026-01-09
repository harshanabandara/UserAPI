package service

import (
	"context"
	"testing"
	"userapi/app/internal/core/domain"
)

type MockUserRepository struct {
	CreateUserFn       func(ctx context.Context, user domain.User) (domain.User, error)
	RetrieveUserFn     func(ctx context.Context, id string) (domain.User, error)
	RetrieveAllUsersFn func(ctx context.Context) ([]domain.User, error)
	UpdateUserFn       func(ctx context.Context, user domain.User, id string) (domain.User, error)
	DeleteUserFn       func(ctx context.Context, id string) error
}

func (m MockUserRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	return m.CreateUserFn(ctx, user)
}

func (m MockUserRepository) RetrieveUser(ctx context.Context, s string) (domain.User, error) {
	return m.RetrieveUserFn(ctx, s)
}

func (m MockUserRepository) RetrieveAllUsers(ctx context.Context) ([]domain.User, error) {
	return m.RetrieveAllUsersFn(ctx)
}

func (m MockUserRepository) UpdateUser(ctx context.Context, s string, user domain.User) (domain.User, error) {
	return m.UpdateUserFn(ctx, user, s)
}

func (m MockUserRepository) DeleteUser(ctx context.Context, s string) error {
	return m.DeleteUserFn(ctx, s)
}

func TestUserServiceImpl_AddUser(t *testing.T) {
	ctx := context.Background()
	user := domain.User{}
	t.Run("Empty User for validation", func(t *testing.T) {
		repo := MockUserRepository{}
		service := NewUserService(repo)
		_, err := service.AddUser(ctx, user)
		if err == nil {
			t.Fatalf("expected validation error")
		}
	})
}
