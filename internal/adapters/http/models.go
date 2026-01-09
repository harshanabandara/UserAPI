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

type UserRequest struct {
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Age       int    `json:"age,omitempty"`
	Status    string `json:"status,omitempty"`
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
