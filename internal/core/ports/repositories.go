package ports

import (
	"context"

	"userapi/app/internal/core/domain"
)

type UserRepository interface {
	CreateUser(context.Context, domain.User) (domain.User, error)
	RetrieveUser(context.Context, string) (domain.User, error)
	RetrieveAllUsers(context.Context) ([]domain.User, error)
	UpdateUser(context.Context, string, domain.User) (domain.User, error)
	DeleteUser(context.Context, string) error
	Close() error
}
