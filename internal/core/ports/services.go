package ports

import (
	"context"

	"userapi/app/internal/core/domain"
)

type UserService interface {
	AddUser(context.Context, domain.User) (domain.User, error)
	GetUserById(context.Context, string) (domain.User, error)
	GetAllUsers(context.Context) ([]domain.User, error)
	UpdateUserByID(context.Context, string, domain.User) (domain.User, error)
	DeleteUserByID(context.Context, string) error
}
