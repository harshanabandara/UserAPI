package main

import (
	"log/slog"

	"userapi/app/internal/adapters/db"
	"userapi/app/internal/adapters/http"
	"userapi/app/internal/adapters/service"
	"userapi/app/internal/core/ports"

	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
)

// @title User Management API
// @version 1.0
// @description This api allow to create, modify,delete, and retrieve user records.

func main() {
	var userRepository ports.UserRepository = db.NewPostgresRepository()
	defer userRepository.Close()
	requestValidator := validator.New()
	var validator ports.Validator = requestValidator
	var userService ports.UserService = service.NewUserService(userRepository, validator)
	server := http.NewServer(userService, validator)
	err := server.Start()
	defer server.Stop()
	if err != nil {
		slog.Error("Could not start the server", "error", err)
		return
	}
}
