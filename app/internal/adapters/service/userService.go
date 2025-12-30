package service

import (
	"errors"
	"userapi/app/internal/core/domain"
	"userapi/app/internal/core/ports"
)

type UserServiceImpl struct {
	UserRepository ports.UserRepository
}

func NewUserService(UserRepository ports.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{UserRepository: UserRepository}
}

func (u *UserServiceImpl) AddUser(user domain.User) (domain.User, error) {
	if user.FirstName == "" || user.LastName == "" || user.Email == "" {
		return user, errors.New("firstName or lastName or email is empty")
	}
	repository := u.UserRepository
	newUser, err := repository.CreateUser(user)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

func (u *UserServiceImpl) GetUserById(userId string) (domain.User, error) {
	if userId == "" {
		return domain.User{}, errors.New("user id is empty")
	}
	user, err := u.UserRepository.RetrieveUser(userId)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *UserServiceImpl) UpdateUserByID(userId string, user domain.User) (domain.User, error) {
	if userId == "" {
		return domain.User{}, errors.New("user id is empty")
	}
	user, err := u.UserRepository.UpdateUser(userId, user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *UserServiceImpl) DeleteUserByID(userId string) error {
	if userId == "" {
		return errors.New("user id is empty")
	}
	err := u.UserRepository.DeleteUser(userId)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserServiceImpl) GetAllUsers() ([]domain.User, error) {
	users, err := u.UserRepository.RetrieveAllUsers()
	if err != nil {
		return make([]domain.User, 0), err
	}
	return users, nil
}
