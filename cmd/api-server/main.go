package main

import (
	"log"
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
	var userService ports.UserService = service.NewUserService(userRepository, requestValidator)
	var server = http.NewServer(userService, requestValidator)
	err := server.Start()
	if err != nil {
		log.Fatal("could not start the server", err)
		return
	}
	defer server.Stop()
}
