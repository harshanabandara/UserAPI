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
	var userRepository ports.UserRepository = db.NewSqlcRepository()
	defer func(userRepository ports.UserRepository) {
		err := userRepository.Close()
		if err != nil {
			log.Fatal("could not close the user repository", err)
		}
	}(userRepository)
	requestValidator := validator.New()
	var userService ports.UserService = service.NewUserService(userRepository, requestValidator)
	var server = http.NewServer(userService, requestValidator)
	err := server.Start()
	if err != nil {
		log.Fatal("could not start the server", err)
		return
	}
	defer func(server *http.Server) {
		err := server.Stop()
		if err != nil {
			log.Fatal("could not stop the server", err)
		}
	}(server)
}
