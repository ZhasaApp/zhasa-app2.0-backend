-- name: CreateBranchDirector :one
INSERT INTO branch_directors (user_id, branch_id)
VALUES ($1, $2) RETURNING id;

-- name: CreateSalesManagerGoalByType :exec
INSERT INTO sales_manager_goals_by_types (sales_manager_id, from_date, to_date, amount, type_id)
VALUES ($1, $2, $3, $4, $5);

-- name: GetBranchDirectorByUserId :many
SELECT *
FROM branch_directors_view bdv
WHERE bdv.user_id = $1;

-- name: GetBranchDirectorByBranchId :one
SELECT *
FROM branch_directors_view bdv
WHERE bdv.branch_id = $1;

-- name: SetSmGoalBySaleType :exec
INSERT INTO sales_manager_goals_by_types (from_date, to_date, amount, sales_manager_id, type_id)
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (from_date, to_date, sales_manager_id, type_id)
DO
UPDATE SET amount = EXCLUDED.amount;

-- name: GetSMGoal :one
SELECT COALESCE(amount, 0)
FROM sales_manager_goals_by_types
WHERE from_date = $1
  AND to_date = $2
  AND sales_manager_id = $3
  AND type_id = $4;

-- name: SetBranchGoalBySaleType :exec
INSERT INTO branch_goals_by_types (from_date, to_date, amount, branch_id, type_id)
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (from_date, to_date, sales_manager_id, type_id)
DO
UPDATE SET amount = EXCLUDED.amount;

-- name: GetBranchGoalBySaleType :one
SELECT COALESCE(amount, 0)
FROM branch_goals_by_types
WHERE from_date = $1
  AND to_date = $2
  AND branch_id = $3
  AND type_id = $4;