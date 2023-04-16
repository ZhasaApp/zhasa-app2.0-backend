-- name: CreateUser :exec
INSERT INTO users (email, password, avatar_url)
VALUES ($1, $2, $3);

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdatePassword :exec
UPDATE users
SET password = $1
WHERE email = $2;
