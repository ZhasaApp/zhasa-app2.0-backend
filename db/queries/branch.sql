-- name: GetBranchesByRating :many
-- Get Ranked Branches
WITH sales_summary AS (SELECT b.id          AS branch_id,
                              SUM(s.amount) AS total_sales_amount
                       FROM sales s
                                INNER JOIN sales_managers sm ON s.sales_manager_id = sm.id
                                INNER JOIN branches b ON sm.branch_id = b.id
                       WHERE s.sale_date BETWEEN $1 AND $2
                       GROUP BY b.id),
     goal_summary AS (SELECT b.id            AS branch_id,
                             SUM(smg.amount) AS total_goal_amount
                      FROM sales_manager_goals smg
                               INNER JOIN sales_managers sm ON smg.sales_manager_id = sm.id
                               INNER JOIN branches b ON sm.branch_id = b.id
                      WHERE smg.from_date = $1
                        AND smg.to_date = $2
                      GROUP BY b.id)
SELECT b.id          AS branch_id,
       b.title       AS branch_title,
       b.branch_key  AS branch_key,
       b.description AS description,
       COALESCE(ss.total_sales_amount / NULLIF(smg.total_goal_amount, 0), 0) ::float AS ratio
FROM branches b
         LEFT JOIN sales_summary ss ON b.id = ss.branch_id
         LEFT JOIN goal_summary smg ON b.id = smg.branch_id
ORDER BY ratio DESC;

-- name: GetBranchById :one
SELECT *
FROM branches
WHERE id = $1;

-- name: CreateBranch :exec
INSERT INTO branches (title, description, branch_key)
VALUES ($1, $2, $3);


-- name: GetBranchSumsByType :many
-- get the sales sums for a specific branch and each sale type within the given period.
SELECT st.id         AS sale_type_id,
       st.title      AS sale_type_title,
       SUM(s.amount) AS total_sales
FROM sale_types st
         JOIN sales s ON st.id = s.sale_type_id
         JOIN sales_managers sm ON s.sales_manager_id = sm.id
WHERE sm.branch_id = $1
  AND s.sale_date BETWEEN $2 AND $3
GROUP BY st.id
ORDER BY st.id ASC;


-- name: GetBranchGoalByGivenDateRange :one
SELECT COALESCE(bg.amount, 0) AS goal_amount
FROM branch_goals bg
WHERE bg.id = $1
  AND bg.from_date = $2
  AND bg.to_date = $3;
