package service

import (
	"context"
	"errors"

	"userapi/app/internal/core/domain"

	"github.com/google/uuid"
)

type MockUserServiceImpl struct {
	users map[string]domain.User
}

func NewMockUserServiceImpl() *MockUserServiceImpl {
	return &MockUserServiceImpl{users: make(map[string]domain.User)}
}

func generateUUID() string {
	return uuid.New().String()
}

func (m MockUserServiceImpl) AddUser(ctx context.Context, user domain.User) (domain.User, error) {
	_ = ctx
	userId := generateUUID()
	user.UserID = userId
	m.users[userId] = user
	return user, nil
}

func (m MockUserServiceImpl) GetUserById(ctx context.Context, s string) (domain.User, error) {
	_ = ctx
	user, ok := m.users[s]
	if !ok {
		return user, errors.New("user not found")
	}
	return user, nil
}

func (m MockUserServiceImpl) UpdateUserByID(ctx context.Context, s string, user domain.User) (domain.User, error) {
	_ = ctx
	currUser, ok := m.users[s]
	if !ok {
		return user, errors.New("user not found")
	}
	if user.FirstName != "" {
		currUser.FirstName = user.FirstName
	}
	if user.LastName != "" {
		currUser.LastName = user.LastName
	}
	if user.Email != "" {
		currUser.Email = user.Email
	}
	if user.Phone != "" {
		currUser.Phone = user.Phone
	}
	if user.Age != 0 {
		currUser.Age = user.Age
	}
	if user.Status == domain.INACTIVE {
		currUser.Status = domain.INACTIVE
	}
	m.users[s] = currUser
	return currUser, nil
}

func (m MockUserServiceImpl) DeleteUserByID(ctx context.Context, s string) error {
	_ = ctx
	delete(m.users, s)
	return nil
}

func (m MockUserServiceImpl) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	_ = ctx
	users := make([]domain.User, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}
