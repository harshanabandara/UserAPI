package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"os"
	"userapi/app/internal/adapters/db"
	"userapi/app/internal/adapters/http"
	"userapi/app/internal/adapters/service"
	"userapi/app/internal/core/ports"
	"userapi/db/sqlc"
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
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer pool.Close()
	queries := sqlc.New(pool)
	var userRepository ports.UserRepository = db.NewSqlcRepository(queries)

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
