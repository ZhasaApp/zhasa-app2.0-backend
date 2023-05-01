-- name: CreateUser :exec
INSERT INTO users (phone, first_name, last_name)
VALUES ($1, $2, $3);

-- name: GetUserByPhone :one
SELECT *
FROM users
WHERE phone = $1;

-- name: CreateUserCode :one
INSERT INTO users_codes(user_id, code)
VALUES ($1, $2)
RETURNING id;

-- name: GetUserCode :one
SELECT *
FROM users_codes
WHERE user_id = $1
ORDER BY created_at DESC LIMIT 1;
