-- name: CreateUser :one
INSERT INTO users (phone, first_name, last_name)
VALUES ($1, $2, $3) ON CONFLICT (phone)
DO
UPDATE SET first_name = EXCLUDED.first_name, last_name = EXCLUDED.last_name
    RETURNING id;

-- name: GetUserByPhone :one
SELECT u.id,
       u.phone,
       u.first_name,
       u.last_name,
       u.avatar_url
FROM user_avatar_view u
WHERE u.phone = $1;

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
INSERT INTO user_brand_sale_type_goals (user_id, brand_id, sale_type_id, value, from_date, to_date)
VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (user_id, brand_id, sale_type_id, from_date, to_date) DO
UPDATE
    SET value = $4;


-- name: GetUsersByBranchBrandRole :many
SELECT u.id,
       u.first_name,
       u.last_name,
       u.avatar_url
FROM user_avatar_view u
         JOIN user_brands ub ON u.id = ub.user_id AND ub.brand_id = $1
         JOIN branch_users bu ON u.id = bu.user_id AND bu.branch_id = $2
         JOIN user_roles ur ON u.id = ur.user_id AND ur.role_id = $3;

-- name: AddBrandToUser :exec
INSERT INTO user_brands (user_id, brand_id)
VALUES ($1, $2) ON CONFLICT DO NOTHING;

-- name: AddRoleToUser :exec
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2) ON CONFLICT DO NOTHING;

-- name: AddUserToBranch :exec
INSERT INTO branch_users (user_id, branch_id)
VALUES ($1, $2) ON CONFLICT DO NOTHING;

-- name: GetUsersWithoutRoles :many
SELECT u.id,
       u.phone,
       u.first_name,
       u.last_name,
       u.created_at
FROM users u
    LEFT JOIN user_roles ur ON u.id = ur.user_id
WHERE ur.user_id IS NULL AND (u.last_name || ' ' || u.first_name) ILIKE @search::text || '%'
ORDER BY u.created_at DESC
LIMIT 25;

-- name: GetUsersWithBranchRolesBrands :many
WITH Counted AS (
    SELECT u.id,
           u.first_name,
           u.last_name,
           b.title                    AS branch_title,
           STRING_AGG(bs.title, ', ') AS brands,
           COUNT(*) OVER()            AS total_count
    FROM users u
             JOIN user_roles ur ON u.id = ur.user_id
             JOIN roles r ON ur.role_id = r.id AND r.key = $1
             JOIN branch_users bu ON u.id = bu.user_id
             JOIN user_brands ub ON u.id = ub.user_id
             JOIN brands bs ON ub.brand_id = bs.id
             JOIN branches b ON bu.branch_id = b.id
    GROUP BY u.id, u.first_name, u.last_name, b.title
)
SELECT id,
       first_name,
       last_name,
       branch_title,
       brands,
       total_count
FROM Counted
ORDER BY first_name, last_name, id DESC
LIMIT $2 OFFSET $3;

-- name: UpdateUser :exec
UPDATE users
SET first_name = $1, last_name = $2, phone = $3
WHERE id = $4;

-- name: UpdateUserBranch :exec
UPDATE branch_users
SET branch_id = $1
WHERE user_id = $2;
