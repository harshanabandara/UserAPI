package http

import "userapi/app/internal/core/domain"

type UserResponse struct {
	UserID    string `json:"userId,omitempty"`
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Age       int    `json:"age,omitempty"`
	Status    string `json:"status,omitempty"`
}

type CreateUserRequest struct {
	FirstName string `json:"firstname" validate:"required,min=2,max=50"`
	LastName  string `json:"lastname" validate:"required,min=2,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone,omitempty" validate:"omitempty,e164"`
	Age       int    `json:"age,omitempty" validate:"omitempty,gte=0,lte=150"`
}

// UserRequest The generic user request
type UserRequest struct {
	FirstName string `json:"firstname,omitempty" validate:"omitempty,min=2,max=50"`
	LastName  string `json:"lastname,omitempty" validate:"omitempty,min=2,max=50"`
	Email     string `json:"email,omitempty" validate:"omitempty,email"`
	Phone     string `json:"phone,omitempty" validate:"omitempty,e164"`
	Age       int    `json:"age,omitempty" validate:"omitempty,gte=0,lte=150"`
	Status    string `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
}

func parseUserToUserDTO(user domain.User) UserResponse {
	return UserResponse{
		UserID:    user.UserID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		Age:       user.Age,
		Status:    user.Status.String(),
	}
}

func (request CreateUserRequest) getUser() domain.User {
	user := domain.User{}
	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.Email = request.Email
	user.Phone = request.Phone
	user.Age = request.Age
	return user
}

func (request UserRequest) getUser() domain.User {
	user := domain.User{}
	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.Email = request.Email
	user.Phone = request.Phone
	user.Age = request.Age
	if request.Status == "active" {
		user.Status = domain.ACTIVE
	} else if request.Status == "inactive" {
		user.Status = domain.INACTIVE
	}
	return user
}
