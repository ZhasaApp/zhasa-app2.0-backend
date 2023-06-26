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
