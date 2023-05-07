-- name: CreateBranchDirector :one
INSERT INTO branch_directors (user_id, branch_id)
VALUES ($1, $2) RETURNING id;

-- name: CreateSalesManagerGoal :exec
INSERT INTO sales_manager_goals (sales_manager_id, from_date, to_date, amount)
VALUES ($1, $2, $3, $4);

-- name: GetBranchDirectorByUserId :one
SELECT * FROM branch_directors_view bdv
WHERE bdv.user_id = $1;