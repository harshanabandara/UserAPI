package main

import (
	"UserApi/internal/adapters/db"
	"UserApi/internal/adapters/http"
	"UserApi/internal/adapters/service"
	"UserApi/internal/core/ports"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func main() {
	host := getEnv("PG_HOST", "localhost")
	port := getEnv("PG_PORT", "5432")
	user := getEnv("PG_USER", "postgres")
	password := getEnv("PG_PASSWORD", "yaalalabs")
	databaseName := getEnv("PG_DATABASE", "userapi")
	sslmode := getEnv("PG_SSLMODE", "disable")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, databaseName, sslmode)
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	var userRepository ports.UserRepository = db.NewSQLUserRepository(database)
	var userService ports.UserService = service.NewUserService(userRepository)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	var server = http.Server{UserService: userService, Router: r}
	err = server.Start()
	if err != nil {
		return
	}

}
