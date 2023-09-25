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

