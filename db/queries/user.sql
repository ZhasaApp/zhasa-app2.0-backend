-- name: CreateUser :exec
INSERT INTO users (phone, first_name, last_name)
VALUES ($1, $2, $3);

-- name: GetUserByPhone :one
SELECT *
FROM users
WHERE phone = $1;

-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1;

-- name: CreateUserCode :one
INSERT INTO users_codes(user_id, code)
VALUES ($1, $2)
RETURNING id;

-- name: GetAuthCodeById :one
SELECT *
FROM users_codes
WHERE id = $1;
