package main

import (
	"UserApi/internal/adapters/http"
	"UserApi/internal/adapters/service"
	"UserApi/internal/core/ports"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	var userService ports.UserService = service.NewMockUserServiceImpl()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	var server = http.Server{UserService: userService, Router: r}
	err := server.Start()
	if err != nil {
		return
	}
}
