package ports

import (
	"UserApi/internal/core/domain"
)

type UserService interface {
	AddUser(user domain.User) (domain.User, error)
	GetUserById(string) (domain.User, error)
	GetAllUsers() ([]domain.User, error)
	UpdateUserByID(string, domain.User) (domain.User, error)
	DeleteUserByID(string) error
}
