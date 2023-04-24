-- name: CreateUser :exec
INSERT INTO users (phone, first_name, last_name)
VALUES ($1, $2, $3);

-- name: GetUserByPhone :one
SELECT * FROM users
WHERE phone = $1;
