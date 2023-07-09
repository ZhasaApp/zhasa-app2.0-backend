-- name: GetBranchById :one
SELECT *
FROM branches
WHERE id = $1;

-- name: CreateBranch :exec
INSERT INTO branches (title, description, branch_key)
VALUES ($1, $2, $3);

-- name: GetBranchGoalByGivenDateRange :one
SELECT COALESCE(bg.amount, 0) AS goal_amount
FROM branch_goals_by_types bg
WHERE bg.id = $1
  AND bg.from_date = $2
  AND bg.to_date = $3
  AND bg.type_id = $4;



