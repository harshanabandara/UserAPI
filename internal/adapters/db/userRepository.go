package db

import (
	"UserApi/internal/core/domain"
	"database/sql"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewSQLUserRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (u *UserRepositoryImpl) RetrieveUser(s string) (domain.User, error) {

	panic("implement me")
}

func (u *UserRepositoryImpl) CreateUser(user domain.User) (domain.User, error) {

	panic("implement me")
}

func (u *UserRepositoryImpl) UpdateUser(userId string, user domain.User) (domain.User, error) {

	panic("implement me")
}

func (u *UserRepositoryImpl) DeleteUser(userId string) error {

	panic("implement me")
}

func (u *UserRepositoryImpl) Init() error {

	panic("implement me")
}

func (u *UserRepositoryImpl) RetrieveAllUsers() ([]domain.User, error) {
	panic("implement me")
}
