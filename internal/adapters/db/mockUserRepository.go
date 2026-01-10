package db

import (
	"context"
	"errors"
	"userapi/app/internal/core/domain"

	"github.com/google/uuid"
)

type MockUserRepository struct {
	users map[string]domain.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{users: make(map[string]domain.User)}
}

func generateUUID() string {
	return uuid.New().String()
}

func (m *MockUserRepository) RetrieveUser(ctx context.Context, s string) (domain.User, error) {
	_ = ctx
	if user, ok := m.users[s]; ok {
		return user, nil
	}
	return domain.User{}, errors.New("user not found")
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	_ = ctx
	userId := generateUUID()
	user.UserID = userId
	m.users[userId] = user
	return user, nil
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, s string, user domain.User) (domain.User, error) {
	//get the user.
	_ = ctx
	currentUser, ok := m.users[s]
	if !ok {
		// Not trying to create a new user.
		return domain.User{}, errors.New("user not found")
	}
	if user.FirstName != "" {
		currentUser.FirstName = user.FirstName
	}
	if user.LastName != "" {
		currentUser.LastName = user.LastName
	}
	if user.Email != "" {
		currentUser.Email = user.Email
	}
	if user.Age != 0 {
		currentUser.Age = user.Age
	}
	if user.Status != 0 {
		currentUser.Status = user.Status
	}
	m.users[s] = currentUser
	return currentUser, nil
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, s string) error {
	_ = ctx
	delete(m.users, s)
	return nil
}

func (m *MockUserRepository) RetrieveAllUsers(ctx context.Context) ([]domain.User, error) {
	_ = ctx
	users := make([]domain.User, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}
