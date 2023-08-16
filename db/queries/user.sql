-- name: CreateUser :one
INSERT INTO users (phone, first_name, last_name)
VALUES ($1, $2, $3) RETURNING id;

-- name: GetUserByPhone :one
SELECT *
FROM user_avatar_view
WHERE phone = $1;

-- name: GetUserById :one
SELECT *
FROM user_avatar_view
WHERE id = $1;

-- name: CreateUserCode :one
INSERT INTO users_codes(user_id, code)
VALUES ($1, $2) RETURNING id;

-- name: GetAuthCodeById :one
SELECT *
FROM users_codes
WHERE id = $1;

-- name: UploadUserAvatar :exec
INSERT INTO users_avatars(user_id, avatar_url)
VALUES ($1, $2) ON CONFLICT (user_id)
DO
UPDATE SET avatar_url = EXCLUDED.avatar_url;

-- name: DeleteUserAvatar :exec
DELETE
FROM users_avatars
WHERE user_id = $1;


