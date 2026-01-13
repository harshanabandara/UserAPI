-- name: RetrieveUserById :one
SELECT * FROM users WHERE user_id = $1 LIMIT 1;

-- name: RetrieveAllUsers :many
SELECT * FROM users;

-- name: CreateUserDefault :one
INSERT INTO users (
    first_name, last_name, email, phone, age
) VALUES (
             $1, $2, $3,$4, $5
          )
RETURNING *;

-- name: CreateUser :one
INSERT INTO users (
    first_name, last_name, email, phone, age, status
) VALUES (
             $1, $2, $3,$4, $5, $6
         )
RETURNING *;

-- name: DeleteUserById :exec
DELETE FROM users WHERE user_id = $1;

-- name: UpdateUserById :one
UPDATE users
SET
    first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name  = COALESCE(sqlc.narg('last_name'), last_name),
    email      = COALESCE(sqlc.narg('email'), email),
    age        = COALESCE(sqlc.narg('age'), age),
    phone      = COALESCE(sqlc.narg('phone'), phone),
    status     = COALESCE(sqlc.narg('status'), status)
WHERE user_id = sqlc.arg('user_id')
RETURNING user_id, first_name, last_name, email, age, phone, status;
