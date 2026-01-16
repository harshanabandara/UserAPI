package service

/*
* This test file contains the test cases to validate the behaviour of the userServiceImpl.
* The UserService is responsible to validate the following data.
*	1. Payload in the incoming requests.
*	2. Data received from the user repository.
 */

import (
	"context"
	"errors"
	"log"
	"testing"
	"userapi/app/internal/core/domain"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type MockUserRepository struct {
	CreateUserFn       func(ctx context.Context, user domain.User) (domain.User, error)
	RetrieveUserFn     func(ctx context.Context, id string) (domain.User, error)
	RetrieveAllUsersFn func(ctx context.Context) ([]domain.User, error)
	UpdateUserFn       func(ctx context.Context, user domain.User, id string) (domain.User, error)
	DeleteUserFn       func(ctx context.Context, id string) error
}

func (m MockUserRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	return m.CreateUserFn(ctx, user)
}

func (m MockUserRepository) RetrieveUser(ctx context.Context, s string) (domain.User, error) {
	return m.RetrieveUserFn(ctx, s)
}

func (m MockUserRepository) RetrieveAllUsers(ctx context.Context) ([]domain.User, error) {
	return m.RetrieveAllUsersFn(ctx)
}

func (m MockUserRepository) UpdateUser(ctx context.Context, s string, user domain.User) (domain.User, error) {
	return m.UpdateUserFn(ctx, user, s)
}

func (m MockUserRepository) DeleteUser(ctx context.Context, s string) error {
	return m.DeleteUserFn(ctx, s)
}

func (m MockUserRepository) Close() error {
	return nil
}

func TestUserServiceImpl_AddUser(t *testing.T) {
	ctx := context.Background()
	entityValidator := validator.New()

	//Add User
	t.Run("Empty User for validation", func(t *testing.T) {
		user := domain.User{}
		repo := MockUserRepository{}
		service := NewUserService(repo, entityValidator)
		_, err := service.AddUser(ctx, user)
		if err == nil {
			t.Fatalf("expected validation error")
		}
	})
	t.Run("no email address", func(t *testing.T) {
		user := domain.User{
			FirstName: "John",
			LastName:  "Doe",
		}
		repo := MockUserRepository{}
		service := NewUserService(repo, entityValidator)
		_, err := service.AddUser(ctx, user)
		if err == nil {
			t.Fatalf("expected validation error for a user without email address")
		}
	})
	t.Run("invalid email address", func(t *testing.T) {
		user := domain.User{
			FirstName: "John",
			LastName:  "Smith",
			Email:     "johnsmith.com",
		}

		repo := MockUserRepository{}
		service := NewUserService(repo, entityValidator)
		_, err := service.AddUser(ctx, user)
		if err == nil {
			t.Fatalf("expected validation error for a user with invalid email address")
		}
	})
	t.Run("firstName too short", func(t *testing.T) {
		user := domain.User{
			FirstName: "J",
			LastName:  "Smith",
			Email:     "johnsmith@mail.com",
		}
		repo := MockUserRepository{}
		service := NewUserService(repo, entityValidator)
		_, err := service.AddUser(ctx, user)
		if err == nil {
			t.Fatalf("expected validation error for a user with invalid first name length")
		}
	})
	t.Run("Age check", func(t *testing.T) {
		user := domain.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@mail.com",
			Age:       -1,
		}
		repo := MockUserRepository{}
		service := NewUserService(repo, entityValidator)
		_, err := service.AddUser(ctx, user)
		if err == nil {
			t.Fatalf("Expected validation error for a user with invalid age check")
		}
	})
	t.Run("Add a valid user test", func(t *testing.T) {
		user := domain.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@mail.com",
			Phone:     "+94700000000",
			Age:       28,
		}
		repo := MockUserRepository{}
		userId := uuid.New().String()
		repo.CreateUserFn = func(ctx context.Context, user domain.User) (domain.User, error) {
			return domain.User{
				UserID:    userId,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@mail.com",
				Phone:     "+94700000000",
				Age:       28,
			}, nil
		}
		service := NewUserService(repo, entityValidator)
		returnUser, err := service.AddUser(ctx, user)
		if err != nil {
			t.Fatal("Failed to add a valid user", err.Error())
		}
		switch {
		case userId != returnUser.UserID:
			t.Fatalf("User ID does not match")
		case user.FirstName != returnUser.FirstName:
			t.Fatalf("First name does not match")
		case user.LastName != returnUser.LastName:
			t.Fatalf("Last name does not match")
		case user.Email != returnUser.Email:
			t.Fatalf("Email does not match")
		case user.Age != returnUser.Age:
			t.Fatalf("Age does not match")
		case user.Phone != returnUser.Phone:
			t.Fatalf("Phone number does not match")
		}
	})
	t.Run("Check repository error", func(t *testing.T) {
		repo := MockUserRepository{}
		repo.CreateUserFn = func(ctx context.Context, user domain.User) (domain.User, error) {
			return domain.User{}, errors.New("mock repository error")
		}
		service := NewUserService(repo, entityValidator)
		user := domain.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@mail.com",
			Phone:     "+94700000000",
			Age:       28,
		}
		_, err := service.AddUser(ctx, user)
		if err == nil {
			t.Fatal("expected repository error, but did not receive via the user service")
		}
	})
	t.Run("Invalid data from the repo", func(t *testing.T) {
		validUser := domain.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@mail.com",
		}
		invalidUser := domain.User{
			UserID:    "invalidUserId",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@mail.com",
		}
		repo := MockUserRepository{}
		repo.CreateUserFn = func(ctx context.Context, user domain.User) (domain.User, error) {
			return invalidUser, nil
		}
		service := NewUserService(repo, entityValidator)
		_, err := service.AddUser(ctx, validUser)
		if err == nil {
			t.Fatal("expected validation error, but did not receive via the user service")
		}
	})

	//Get User
	t.Run("Get user with empty user id", func(t *testing.T) {
		repo := MockUserRepository{}
		service := NewUserService(repo, entityValidator)
		_, err := service.GetUserById(ctx, "")
		if err == nil {
			t.Fatal("expected error, but did not receive via the user service")
		}
	})
	t.Run("Get user with nonexistent user id", func(t *testing.T) {
		repo := MockUserRepository{}
		repo.RetrieveUserFn = func(ctx context.Context, id string) (domain.User, error) {
			return domain.User{}, errors.New("user id not found error")
		}
		service := NewUserService(repo, entityValidator)
		userId := "user1"
		_, err := service.GetUserById(ctx, userId)
		if err == nil {
			t.Fatal("expected error, but did not receive via the user service")
		}
	})
	t.Run("Get user with valid user id", func(t *testing.T) {
		repo := MockUserRepository{}
		userId := uuid.New().String()
		repo.RetrieveUserFn = func(ctx context.Context, id string) (domain.User, error) {
			return domain.User{
				UserID:    userId,
				FirstName: "John",
				LastName:  "Smith",
				Email:     "john.smith@mail.com",
			}, nil
		}
		service := NewUserService(repo, entityValidator)
		user, err := service.GetUserById(ctx, userId)
		if err != nil {
			t.Fatal("Unexpected error while get user", err.Error())
		}
		if user.UserID != userId {
			t.Fatalf("User ID does not match")
		}
	})
	t.Run("Get user with a different user id", func(t *testing.T) {
		repo := MockUserRepository{}
		userId := uuid.New().String()
		repo.RetrieveUserFn = func(ctx context.Context, id string) (domain.User, error) {
			return domain.User{
				UserID:    uuid.New().String(),
				FirstName: "John",
				LastName:  "Smith",
				Email:     "john.smith@mail.com",
			}, nil
		}
		service := NewUserService(repo, entityValidator)
		_, err := service.GetUserById(ctx, userId)
		if err == nil {
			t.Fatal("Error expected. Should validate the return data from the DB", err.Error())
		}
	})
	t.Run("Get user with invalid user id", func(t *testing.T) {
		repo := MockUserRepository{}
		userId := "invalidUserId"
		userService := NewUserService(repo, entityValidator)
		_, err := userService.GetUserById(ctx, userId)
		if err == nil {
			t.Fatal("Error expected. Should validate the user id", err.Error())
		}
	})
	t.Run("Get user invalid data from database", func(t *testing.T) {
		repo := MockUserRepository{}
		userId := uuid.New().String()
		repo.RetrieveUserFn = func(ctx context.Context, id string) (domain.User, error) {
			return domain.User{
				UserID:    userId,
				FirstName: "John",
				LastName:  "Smith",
				Email:     "john.com",
			}, nil
		}
		service := NewUserService(repo, entityValidator)
		_, err := service.GetUserById(ctx, userId)
		if err == nil {
			t.Fatal("Error expected. Should validate the user id", err.Error())
		}
	})
	t.Run("Get user database error", func(t *testing.T) {
		repo := MockUserRepository{}
		userId := uuid.New().String()
		repo.RetrieveUserFn = func(ctx context.Context, id string) (domain.User, error) {
			return domain.User{}, errors.New("mock db error")
		}
		service := NewUserService(repo, entityValidator)
		_, err := service.GetUserById(ctx, userId)
		if err == nil {
			t.Fatal("Error expected. User service should forward the error", err)
		}
	})

	//Get All Users
	t.Run("Get users database error", func(t *testing.T) {
		repo := MockUserRepository{}
		repo.RetrieveAllUsersFn = func(ctx context.Context) ([]domain.User, error) {
			return nil, errors.New("mock db error")
		}
		userService := NewUserService(repo, entityValidator)
		_, err := userService.GetAllUsers(ctx)
		if err == nil {
			t.Fatal("Database Error expected. User service should forward the error")
		}
	})
	t.Run("Get users return invalid data", func(t *testing.T) {
		users := make([]domain.User, 2)
		users = append(users, domain.User{
			UserID:    uuid.New().String(),
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@mailcom",
		},
			domain.User{
				UserID:    uuid.New().String(),
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "janedoe@mail.com",
				Age:       27,
			})

		repo := MockUserRepository{}
		repo.RetrieveAllUsersFn = func(ctx context.Context) ([]domain.User, error) {
			return users, nil
		}
		userService := NewUserService(repo, entityValidator)
		usrList, err := userService.GetAllUsers(context.Background())
		if err == nil {
			log.Fatal("User service should not return invalid data", usrList)
		}
	})
	t.Run("Get users return correct data", func(t *testing.T) {
		users := []domain.User{
			{
				UserID:    uuid.New().String(),
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@mail.com",
			},
			{
				UserID:    uuid.New().String(),
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane.doe@mail.com",
			},
		}
		repo := MockUserRepository{}
		repo.RetrieveAllUsersFn = func(ctx context.Context) ([]domain.User, error) {
			return users, nil
		}
		userService := NewUserService(repo, entityValidator)
		usrList, err := userService.GetAllUsers(context.Background())
		if err != nil {
			log.Fatal("User service should not return an error for valid data", usrList)
		}
	})

	//Update User
	t.Run("Update user with empty user id", func(t *testing.T) {
		repo := MockUserRepository{}
		userId := ""
		userService := NewUserService(repo, entityValidator)
		_, err := userService.UpdateUserByID(ctx, userId, domain.User{
			FirstName: "Jane",
		})
		if err == nil {
			t.Fatal("Error expected. Should validate the uuid")
		}
	})
	t.Run("Update user with invalid user id", func(t *testing.T) {
		repo := MockUserRepository{}
		userId := "invalidUserId"
		userService := NewUserService(repo, entityValidator)
		_, err := userService.UpdateUserByID(ctx, userId, domain.User{
			FirstName: "Jane",
		})
		if err == nil {
			t.Fatal("Error expected. Should validate the user id")
		}
	})
	t.Run("Update user database error", func(t *testing.T) {
		repo := MockUserRepository{}
		userId := uuid.New().String()
		repo.UpdateUserFn = func(ctx context.Context, user domain.User, userId string) (domain.User, error) {
			return domain.User{}, errors.New("mock db error")
		}
		userService := NewUserService(repo, entityValidator)
		_, err := userService.UpdateUserByID(ctx, userId, domain.User{
			FirstName: "Jane",
		})
		if err == nil {
			t.Fatal("Error expected. Should return the database error")
		}
	})
	t.Run("Update a valid user entry", func(t *testing.T) {
		repo := MockUserRepository{}
		userId := uuid.New().String()
		repo.UpdateUserFn = func(ctx context.Context, user domain.User, userId string) (domain.User, error) {
			return domain.User{
				UserID:    userId,
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane.doe@mail.com",
				Age:       27,
			}, nil
		}
		userService := NewUserService(repo, entityValidator)
		_, err := userService.UpdateUserByID(ctx, userId, domain.User{
			Age: 27,
		})
		if err != nil {
			t.Fatal("Unexpected error. Should be able to update the user.", err)
		}
	})

	//Delete user
	t.Run("Delete user database error", func(t *testing.T) {
		repo := MockUserRepository{}
		userId := uuid.New().String()
		repo.DeleteUserFn = func(ctx context.Context, userId string) error {
			return errors.New("mock db error")
		}
		userService := NewUserService(repo, entityValidator)
		err := userService.DeleteUserByID(ctx, userId)
		if err == nil {
			t.Fatal("Error expected. Should return the database error")
		}
	})
	t.Run("Delete user invalid uuid", func(t *testing.T) {
		repo := MockUserRepository{}
		userId := "invalidUserId"
		userService := NewUserService(repo, entityValidator)
		err := userService.DeleteUserByID(ctx, userId)
		if err == nil {
			t.Fatal("Error expected. Should return the database error")
		}
	})

	t.Run("Delete user with valid uuid", func(t *testing.T) {
		repo := MockUserRepository{}
		repo.DeleteUserFn = func(ctx context.Context, userId string) error {
			return nil
		}
		userId := uuid.New().String()
		userService := NewUserService(repo, entityValidator)
		err := userService.DeleteUserByID(ctx, userId)
		if err != nil {
			t.Fatal("Unexpected error. Should be able to delete the user", err)
		}
	})
}
