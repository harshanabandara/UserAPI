package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"userapi/app/internal/adapters/db/user"
	"userapi/app/internal/core/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	q    *sqlc.Queries
	pool *pgxpool.Pool
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func NewPostgresRepository() *PostgresRepository {
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
		return nil
	}
	queries := sqlc.New(pool)
	return &PostgresRepository{q: queries, pool: pool}
}

func (repository *PostgresRepository) Close() error {
	repository.pool.Close()
	return nil
}

func (repository *PostgresRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	params := parseUserToCreateUserParams(user)
	newUser, err := repository.q.CreateUser(ctx, params)
	if err != nil {
		return domain.User{}, err
	}
	user.UserID = newUser.UserID.String()
	user = getUserFromUserRecord(newUser)

	return user, nil
}

func (repository *PostgresRepository) RetrieveUser(ctx context.Context, userId string) (domain.User, error) {
	userUuid, err := uuid.Parse(userId)
	if err != nil {
		return domain.User{}, err
	}
	user, err := repository.q.RetrieveUserById(ctx, userUuid)
	if err != nil {
		return domain.User{}, err
	}
	returnUser := getUserFromUserRecord(user)
	return returnUser, nil
}

func (repository *PostgresRepository) RetrieveAllUsers(ctx context.Context) ([]domain.User, error) {
	allUsers, err := repository.q.RetrieveAllUsers(ctx)
	if err != nil {
		return []domain.User{}, err
	}
	users := make([]domain.User, len(allUsers))
	for index, user := range allUsers {
		users[index] = getUserFromUserRecord(user)
	}
	return users, nil
}

func (repository *PostgresRepository) UpdateUser(ctx context.Context, userId string, user domain.User) (domain.User, error) {
	userUuid, err := uuid.Parse(userId)
	if err != nil {
		return domain.User{}, err
	}
	params := sqlc.UpdateUserByIdParams{}
	params.UserID = userUuid
	if user.FirstName != "" {
		params.FirstName = pgtype.Text{String: user.FirstName, Valid: true}
	} else {
		params.FirstName = pgtype.Text{String: "", Valid: false}
	}
	if user.LastName != "" {
		params.LastName = pgtype.Text{String: user.LastName, Valid: true}
	} else {
		params.LastName = pgtype.Text{String: "", Valid: false}
	}
	if user.Email != "" {
		params.Email = pgtype.Text{String: user.Email, Valid: true}
	} else {
		params.Email = pgtype.Text{String: "", Valid: false}
	}
	if user.Phone != "" {
		params.Phone = pgtype.Text{String: user.Phone, Valid: true}
	} else {
		params.Phone = pgtype.Text{String: "", Valid: false}
	}
	if user.Age != 0 {
		params.Age = pgtype.Int4{
			Int32: int32(user.Age), //nolint:gosec
			Valid: true,
		}
	} else {
		params.Age = pgtype.Int4{Int32: 0, Valid: false}
	}
	switch user.Status {
	case domain.ACTIVE:
		params.Status = sqlc.NullUserStatus{
			UserStatus: sqlc.UserStatusACTIVE,
			Valid:      true,
		}
	case domain.INACTIVE:
		params.Status = sqlc.NullUserStatus{
			UserStatus: sqlc.UserStatusINACTIVE,
			Valid:      true,
		}
	default:
		params.Status = sqlc.NullUserStatus{
			UserStatus: sqlc.UserStatusACTIVE,
			Valid:      false,
		}
	}
	row, err := repository.q.UpdateUserById(ctx, params)
	if err != nil {
		return domain.User{}, err
	}

	returnUser := domain.User{
		UserID:    row.UserID.String(),
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Email:     row.Email,
		Age:       int(row.Age.Int32),
		Status:    getUserStatusFromStatusRecord(row.Status),
	}
	return returnUser, nil
}

func (repository *PostgresRepository) DeleteUser(ctx context.Context, userId string) error {
	userUuid, err := uuid.Parse(userId)
	if err != nil {
		return err
	}
	err = repository.q.DeleteUserById(ctx, userUuid)
	if err != nil {
		return err
	}
	return nil
}

func getStringFromTextRecord(s pgtype.Text) string {
	if s.Valid {
		return s.String
	}
	return ""
}

func getUserStatusFromStatusRecord(s sqlc.NullUserStatus) domain.UserStatus {
	if !s.Valid {
		return domain.UserStatus(0)
	}
	if s.UserStatus == sqlc.UserStatusACTIVE {
		return domain.ACTIVE
	}
	if s.UserStatus == sqlc.UserStatusINACTIVE {
		return domain.INACTIVE
	}
	return domain.ACTIVE
}

func getUserFromUserRecord(userRecord sqlc.User) domain.User {
	user := domain.User{}
	user.FirstName = userRecord.FirstName
	user.LastName = userRecord.LastName
	user.Email = userRecord.Email
	user.Phone = getStringFromTextRecord(userRecord.Phone)
	user.Status = getUserStatusFromStatusRecord(userRecord.Status)
	user.UserID = userRecord.UserID.String()
	user.Age = int(userRecord.Age.Int32)
	return user
}

func parseUserToCreateUserParams(user domain.User) sqlc.CreateUserParams {
	params := sqlc.CreateUserParams{}
	if user.FirstName != "" {
		params.FirstName = user.FirstName
	}
	if user.LastName != "" {
		params.LastName = user.LastName
	}
	if user.Email != "" {
		params.Email = user.Email
	}
	if user.Phone != "" {
		params.Phone = pgtype.Text{
			String: user.Phone,
			Valid:  true,
		}
	} else {
		params.Phone = pgtype.Text{String: "", Valid: false}
	}
	if user.Age != 0 {
		params.Age = pgtype.Int4{
			Int32: int32(user.Age), //nolint:gosec    //check endian impact.
			Valid: true,
		}
	} else {
		params.Age = pgtype.Int4{Int32: 0, Valid: false}
	}
	if user.Status == domain.INACTIVE {
		params.Status = sqlc.NullUserStatus{
			UserStatus: sqlc.UserStatusINACTIVE,
			Valid:      true,
		}
	} else {
		params.Status = sqlc.NullUserStatus{
			UserStatus: sqlc.UserStatusACTIVE,
			Valid:      true,
		}
	}
	return params
}
