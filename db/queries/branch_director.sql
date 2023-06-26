-- name: CreateBranchDirector :one
INSERT INTO branch_directors (user_id, branch_id)
VALUES ($1, $2) RETURNING id;

-- name: CreateSalesManagerGoalByType :exec
INSERT INTO sales_manager_goals_by_types (sales_manager_id, from_date, to_date, amount, type_id)
VALUES ($1, $2, $3, $4, $5);

-- name: GetBranchDirectorByUserId :one
SELECT * FROM branch_directors_view bdv
WHERE bdv.user_id = $1;