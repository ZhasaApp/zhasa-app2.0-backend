-- name: CreateSalesManager :exec
INSERT INTO sales_managers (user_id, branch_id)
VALUES ($1, $2);

-- name: GetSalesManagerSumsByType :one
-- get the sales sums for a specific sales manager and each sale type within the given period.
SELECT st.id         AS sale_type_id,
       st.title      AS sale_type_title,
       SUM(s.amount) AS total_sales
FROM sale_types st
         JOIN sales s ON st.id = s.sale_type_id AND s.sales_manager_id = $1 AND s.sale_date BETWEEN $2 AND $3
WHERE st.id = $4
GROUP BY st.id
ORDER BY st.id ASC;


-- name: AddSaleOrReplace :one
-- add sale into sales by given sale_type_id, amount, date, sales_manager_id and on conflict replace
INSERT INTO sales (sales_manager_id, sale_date, amount, sale_type_id, description)
VALUES ($1, $2, $3, $4, $5) RETURNING *;;

-- name: GetSalesByDate :many
SELECT *
from sales s
WHERE s.sale_date = $1;

-- name: GetSalesManagerByUserId :one
SELECT *
from sales_managers_view s
WHERE s.user_id = $1;

-- name: GetSalesManagerGoalByGivenDateRangeAndSaleType :one
SELECT COALESCE(sg.amount, 0) AS goal_amount
FROM sales_manager_goals_by_types sg
WHERE sg.sales_manager_id = $1
  AND sg.from_date = $2
  AND sg.to_date = $3
  AND sg.type_id = $4;

-- name: GetManagerSales :many
SELECT id, sale_type_id, description, sale_date, amount
FROM sales s
WHERE s.sales_manager_id = $1
ORDER BY s.sale_date DESC LIMIT $2
OFFSET $3;


-- name: GetManagerSalesByPeriod :many
SELECT id, sale_type_id, description, sale_date, amount
FROM sales s
WHERE s.sales_manager_id = $1
  AND s.sale_date BETWEEN $2 AND $3
ORDER BY s.sale_date DESC LIMIT $4
OFFSET $5;

-- name: GetSalesCount :one
SELECT COUNT(*)
FROM sales
WHERE sales_manager_id = $1;

-- name: SetSMRatio :exec
INSERT INTO sales_manager_goals_ratio_by_period
    (from_date, to_date, ratio, sales_manager_id)
VALUES ($1, $2, $3, $4) ON CONFLICT (from_date, to_date, sales_manager_id)
DO
UPDATE SET ratio = EXCLUDED.ratio;


-- name: GetSMRatio :one
SELECT ratio
FROM sales_manager_goals_ratio_by_period smgr
WHERE smgr.from_date = $1
  AND smgr.to_date = $2
  AND smgr.sales_manager_id = $3;