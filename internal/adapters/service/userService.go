package service

import (
	"context"
	"errors"
	"fmt"
	"userapi/app/internal/core/domain"
	"userapi/app/internal/core/ports"
)

type UserServiceImpl struct {
	UserRepository ports.UserRepository
}

func NewUserService(UserRepository ports.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{UserRepository: UserRepository}
}

func (u *UserServiceImpl) AddUser(ctx context.Context, user domain.User) (domain.User, error) {
	if user.FirstName == "" || user.LastName == "" || user.Email == "" {
		return user, errors.New("firstName or lastName or email is empty")
	}
	repository := u.UserRepository
	newUser, err := repository.CreateUser(ctx, user)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

func (u *UserServiceImpl) GetUserById(ctx context.Context, userId string) (domain.User, error) {
	if userId == "" {
		return domain.User{}, errors.New("user id is empty")
	}
	user, err := u.UserRepository.RetrieveUser(ctx, userId)
	if err != nil {
		return domain.User{}, fmt.Errorf("user with id %s not found : %w", userId, err)
	}
	return user, nil
}

func (u *UserServiceImpl) UpdateUserByID(ctx context.Context, userId string, user domain.User) (domain.User, error) {
	if userId == "" {
		return domain.User{}, errors.New("user id is empty")
	}
	user, err := u.UserRepository.UpdateUser(ctx, userId, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("could not update the user with id %s : %w", userId, err)
	}
	return user, nil
}

func (u *UserServiceImpl) DeleteUserByID(ctx context.Context, userId string) error {
	if userId == "" {
		return errors.New("user id is empty")
	}
	err := u.UserRepository.DeleteUser(ctx, userId)
	if err != nil {
		return fmt.Errorf("could not delete the user with id %s : %w", userId, err)
	}
	return nil
}

func (u *UserServiceImpl) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	users, err := u.UserRepository.RetrieveAllUsers(ctx)
	if err != nil {
		return make([]domain.User, 0), fmt.Errorf("could not retrieve the users:  %w", err)
	}
	return users, nil
}
