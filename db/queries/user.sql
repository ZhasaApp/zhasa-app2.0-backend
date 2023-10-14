-- name: CreateUser :one
INSERT INTO users (phone, first_name, last_name)
VALUES ($1, $2, $3) ON CONFLICT (phone)
DO
UPDATE SET first_name = EXCLUDED.first_name, last_name = EXCLUDED.last_name
    RETURNING id;

-- name: GetUserByPhone :one
SELECT *
FROM user_avatar_view
WHERE phone = $1;

-- name: GetUserById :one
SELECT *
FROM user_avatar_view u
         JOIN user_roles ur on u.id = ur.user_id
         JOIN roles r on ur.role_id = r.id
WHERE u.id = $1;

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

-- name: GetUserBranch :one
SELECT b.title, b.id
FROM users u
         JOIN
     branch_users bu ON u.id = bu.user_id
         JOIN branches b ON bu.branch_id = b.id
WHERE u.id = $1;

-- name: SetUserBrandGoal :exec
INSERT INTO user_brand_sale_type_goals (user_brand, sale_type_id, value, from_date, to_date)
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (user_brand, sale_type_id, from_date, to_date) DO
UPDATE
    SET value = $3;
