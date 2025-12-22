package db

import (
	"UserApi/internal/core/domain"
	"errors"
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

func (m *MockUserRepository) RetrieveUser(s string) (domain.User, error) {
	if user, ok := m.users[s]; ok {
		return user, nil
	}
	return domain.User{}, errors.New("user not found")
}

func (m *MockUserRepository) CreateUser(user domain.User) (domain.User, error) {
	userId := generateUUID()
	user.UserID = userId
	m.users[userId] = user
	return user, nil
}

func (m *MockUserRepository) UpdateUser(s string, user domain.User) (domain.User, error) {
	//get the user.
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

func (m *MockUserRepository) DeleteUser(s string) error {
	if _, ok := m.users[s]; ok {
		delete(m.users, s)
	}
	return nil
}

func (m *MockUserRepository) RetrieveAllUsers() ([]domain.User, error) {
	users := make([]domain.User, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserRepository) Init() error {
	return nil
}
