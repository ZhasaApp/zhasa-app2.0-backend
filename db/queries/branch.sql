-- name: GetBranchById :one
SELECT *
FROM branches
WHERE id = $1;

-- name: CreateBranch :exec
INSERT INTO branches (title, description)
VALUES ($1, $2);

-- name: GetBranchBrandGoalByGivenDateRange :one
SELECT COALESCE(bg.value, 0) AS goal_amount
FROM branch_brand_sale_type_goals bg
WHERE bg.branch_brand = $1
  AND bg.from_date = $2
  AND bg.to_date = $3
  AND bg.sale_type_id = $4;

-- name: GetBranches :many
SELECT *
FROM branches;

-- name: GetBranchesByBrandId :many
SELECT b.id, b.title, b.description
FROM branches b
         JOIN branch_brands bb ON b.id = bb.branch_id
WHERE bb.brand_id = $1;

-- name: SetBranchBrandGoal :exec
INSERT INTO branch_brand_sale_type_goals (branch_brand, sale_type_id, value, from_date, to_date)
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (branch_brand, sale_type_id, from_date, to_date) DO
UPDATE
    SET value = $3;

-- name: GetBranchBrandUserByRole :many
SELECT u.id,
       u.first_name,
       u.last_name,
       u.avatar_url,
       b.title AS branch_title,
       b.id    AS branch_id
FROM user_avatar_view u
         JOIN user_brands ub ON u.id = ub.user_id AND ub.brand_id = $1
         JOIN branch_users bu ON u.id = bu.user_id AND bu.branch_id = $2
         JOIN branches b ON bu.branch_id = b.id
         JOIN user_roles ur ON u.id = ur.user_id AND ur.role_id = $3;