-- name: RetrieveUserById :one
SELECT * FROM users WHERE user_id = $1 LIMIT 1;

-- name: RetrieveAllUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (
    first_name, last_name, email, phone, age, status
) VALUES (
             $1, $2, $3,$4, $5, $6
         )
RETURNING *;

-- name: UpdateUserById :exec

UPDATE users
SET
    first_name = COALESCE($1, first_name),
    last_name  = COALESCE($2, last_name),
    email      = COALESCE($3, email),
    age        = COALESCE($4, age),
    phone      = COALESCE($5, phone),
    status     = COALESCE($6, status)
WHERE user_id = $7
RETURNING *;

-- name: DeleteUserById :exec
DELETE FROM users WHERE user_id = $1;
