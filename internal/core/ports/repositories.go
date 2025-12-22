package ports

import "UserApi/internal/core/domain"

type UserRepository interface {
	Init() error
	CreateUser(domain.User) (domain.User, error)
	RetrieveUser(string) (domain.User, error)
	RetrieveAllUsers() ([]domain.User, error)
	UpdateUser(string, domain.User) (domain.User, error)
	DeleteUser(string) error
}
