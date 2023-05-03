-- name: GetBranchesByRating :many
-- Get Ranked Branches
WITH sales_summary AS (
    SELECT
        b.id AS branch_id,
        SUM(s.amount) AS total_sales_amount
    FROM
        sales s
            INNER JOIN sales_managers sm ON s.sales_manager_id = sm.id
            INNER JOIN branches b ON sm.branch_id = b.id
    WHERE
        s.sale_date BETWEEN $1 AND $2
    GROUP BY
        b.id
),
     goal_summary AS (
         SELECT
             b.id AS branch_id,
             SUM(smg.amount) AS total_goal_amount
         FROM
             sales_manager_goals smg
                 INNER JOIN sales_managers sm ON smg.sales_manager_id = sm.id
                 INNER JOIN branches b ON sm.branch_id = b.id
         WHERE
                 smg.from_date = $1
           AND smg.to_date = $2
         GROUP BY
             b.id
     )
SELECT
    b.id AS branch_id,
    b.title AS branch_title,
    b.branch_key AS branch_key,
    b.description AS description,
    COALESCE(ss.total_sales_amount / NULLIF(smg.total_goal_amount, 0), 0)::float AS ratio
FROM
    branches b
        LEFT JOIN sales_summary ss ON b.id = ss.branch_id
        LEFT JOIN goal_summary smg ON b.id = smg.branch_id
ORDER BY
    ratio DESC;

-- name: GetBranchById :one
SELECT * FROM branches
WHERE id = $1;

-- name: CreateBranch :exec
INSERT INTO branches (title, description, branch_key)
VALUES ($1, $2, $3);