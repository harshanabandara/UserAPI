package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"userapi/app/internal/core/domain"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewSQLUserRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (u *UserRepositoryImpl) RetrieveUser(userID string) (domain.User, error) {
	query := "SELECT user_id, first_name, last_name, email, age, status FROM users WHERE user_id = $1"

	var user domain.User
	var userStatus string
	var active = domain.ACTIVE
	err := u.db.QueryRow(query, userID).Scan(
		&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Age, &userStatus,
	)
	switch userStatus {
	case active.String():
		user.Status = domain.ACTIVE
	default:
		user.Status = domain.INACTIVE
	}
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *UserRepositoryImpl) CreateUser(user domain.User) (domain.User, error) {
	query := "INSERT INTO users (first_name, last_name, email, age, phone) VALUES ($1, $2, $3, $4, $5) RETURNING user_id"
	var id string

	err := u.db.QueryRow(query, user.FirstName, user.LastName, user.Email, user.Age, user.Phone).Scan(&id)
	if err != nil {
		return domain.User{}, err
	}
	user.UserID = id
	return user, nil
}

func (u *UserRepositoryImpl) UpdateUser(userId string, user domain.User) (domain.User, error) {
	if userId == "" {
		return domain.User{}, errors.New("user id is required")
	}
	fields := []string{}
	vals := []interface{}{}
	if user.FirstName != "" {
		fields = append(fields, fmt.Sprintf("first_name = $%d", len(fields)+1))
		vals = append(vals, user.FirstName)
	}
	if user.LastName != "" {
		fields = append(fields, fmt.Sprintf("last_name = $%d", len(fields)+1))
		vals = append(vals, user.LastName)
	}
	if user.Email != "" {
		fields = append(fields, fmt.Sprintf("email = $%d", len(fields)+1))
		vals = append(vals, user.Email)
	}
	if user.Age != 0 {
		fields = append(fields, fmt.Sprintf("age = $%d", len(fields)+1))
		vals = append(vals, user.Age)
	}
	if user.Phone != "" {
		fields = append(fields, fmt.Sprintf("phone = $%d", len(fields)+1))
		vals = append(vals, user.Phone)
	}
	if user.Status != 0 {
		fields = append(fields, fmt.Sprintf("status = $%d", len(fields)+1))
		vals = append(vals, user.Status.String())
	}
	vals = append(vals, userId)
	// build the string
	updatedUser := domain.User{}
	var status string
	var inactive = domain.INACTIVE
	query := fmt.Sprintf("UPDATE users SET %s WHERE user_id = $%d RETURNING user_id,first_name, last_name, email, age, phone, status", strings.Join(fields, ","), len(fields)+1)
	err := u.db.QueryRow(query, vals...).Scan(
		&updatedUser.UserID, &updatedUser.FirstName, &updatedUser.LastName, &updatedUser.Email, &updatedUser.Age, &updatedUser.Phone, &status)
	if err != nil {
		return domain.User{}, err
	}
	switch status {
	case inactive.String():
		updatedUser.Status = domain.INACTIVE
	default:
		updatedUser.Status = domain.ACTIVE
	}
	return updatedUser, nil

}

func (u *UserRepositoryImpl) DeleteUser(userId string) error {
	query := "DELETE FROM users WHERE user_id = $1"
	_, err := u.db.Exec(query, userId)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepositoryImpl) Init() error {
	return nil
}

func (u *UserRepositoryImpl) RetrieveAllUsers() ([]domain.User, error) {
	query := "SELECT user_id, first_name, last_name, email, age, status FROM users"
	rows, err := u.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	var inactive = domain.INACTIVE
	for rows.Next() {
		var u domain.User
		var userStatus string
		if err := rows.Scan(&u.UserID, &u.FirstName, &u.LastName, &u.Email, &u.Age, &userStatus); err != nil {
			return nil, err
		}

		switch userStatus {
		case inactive.String():
			u.Status = domain.INACTIVE
		default:
			u.Status = domain.ACTIVE
		}
		users = append(users, u)
	}
	return users, nil
}
