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
WHERE bg.branch_id = $1
  AND bg.from_date = $2
  AND bg.to_date = $3
  AND bg.type_id = $4;


-- name: GetOrderedBranchesByGivenPeriod :many
SELECT b.title,
       b.id,
       b.description,
       COALESCE(r.ratio, 0.0) AS ratio
FROM branches b
         LEFT JOIN
     branches_goals_ratio_by_period r ON b.id = r.branch_id
         AND r.from_date >= $1 AND r.to_date <= $2
ORDER BY ratio DESC LIMIT $3
OFFSET $4;



