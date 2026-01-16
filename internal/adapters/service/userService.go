package service

import (
	"context"
	"errors"
	"fmt"
	"userapi/app/internal/core/domain"
	"userapi/app/internal/core/ports"

	"github.com/go-playground/validator/v10"
)

type UserServiceImpl struct {
	UserRepository ports.UserRepository
	Validator      *validator.Validate
}

func NewUserService(UserRepository ports.UserRepository, validator *validator.Validate) *UserServiceImpl {
	return &UserServiceImpl{UserRepository: UserRepository, Validator: validator}
}

func (u *UserServiceImpl) AddUser(ctx context.Context, user domain.User) (domain.User, error) {
	validationErr := u.Validator.Struct(user)
	if validationErr != nil {
		return domain.User{}, fmt.Errorf("could not add the user. %w", validationErr)
	}
	if user.FirstName == "" || user.LastName == "" || user.Email == "" {
		return user, errors.New("firstName or lastName or email is empty")
	}
	repository := u.UserRepository
	newUser, err := repository.CreateUser(ctx, user)
	if err != nil {
		return newUser, err
	}
	validationErr = u.Validator.Struct(newUser)
	if validationErr != nil {
		return newUser, fmt.Errorf("could not add the user. %w", validationErr)
	}
	return newUser, nil
}

func (u *UserServiceImpl) GetUserById(ctx context.Context, userId string) (domain.User, error) {
	if userId == "" {
		return domain.User{}, errors.New("user id is empty")
	}
	if uuidErr := u.Validator.Var(userId, "uuid"); uuidErr != nil {
		return domain.User{}, errors.New("user id is not valid")
	}
	user, err := u.UserRepository.RetrieveUser(ctx, userId)
	if err != nil {
		return domain.User{}, fmt.Errorf("user with id %s not found : %w", userId, err)
	}
	validationErr := u.Validator.Struct(user)
	if validationErr != nil {
		return domain.User{}, fmt.Errorf("could not utils the retrieved user. %w", validationErr)
	}
	if userId != user.UserID {
		return domain.User{}, errors.New(fmt.Sprintf("invalid user id returned. expected: %s returned %s", userId, user.UserID))
	}
	return user, nil
}

func (u *UserServiceImpl) UpdateUserByID(ctx context.Context, userId string, user domain.User) (domain.User, error) {
	uuidErr := u.Validator.Var(userId, "uuid")
	if uuidErr != nil {
		return domain.User{}, errors.New("user id is not valid")
	}
	validationErr := u.Validator.Struct(user)
	if validationErr != nil {
		return domain.User{}, fmt.Errorf("could not add the user. %w", validationErr)
	}
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
	uuidErr := u.Validator.Var(userId, "uuid")
	if uuidErr != nil {
		return errors.New("user id is not valid")
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
	validationErr := u.Validator.Var(users, "omitempty,dive")
	if validationErr != nil {
		return make([]domain.User, 0), fmt.Errorf("could not utils the retrieved users. %w", validationErr)
	}
	return users, nil
}
