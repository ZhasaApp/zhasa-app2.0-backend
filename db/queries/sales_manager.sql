-- name: CreateSalesManager :exec
INSERT INTO sales_managers (user_id, branch_id)
VALUES ($1, $2);

-- name: GetSalesManagerSumsByType :many
-- get the sales sums for a specific sales manager and each sale type within the given period.
SELECT st.id         AS sale_type_id,
       st.title      AS sale_type_title,
       SUM(s.amount) AS total_sales
FROM sale_types st
         JOIN sales s ON st.id = s.sale_type_id AND s.sales_manager_id = $1 AND s.sale_date BETWEEN $2 AND $3
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

-- name: GetSalesManagerGoalByGivenDateRange :one
SELECT COALESCE(sg.amount, 0) AS goal_amount
FROM sales_manager_goals sg
WHERE sg.sales_manager_id = $1
  AND sg.from_date = $2
  AND sg.to_date = $3;

-- name: GetManagerSales :many
SELECT id, sale_type_id, description, sale_date, amount
FROM sales s
WHERE s.sales_manager_id = $1
ORDER BY s.sale_date DESC LIMIT $2
OFFSET $3;


-- name: GetManagerSalesByPeriod :many
SELECT id, sale_type_id, description, sale_date, amount
FROM sales s
WHERE s.sales_manager_id = $1 AND s.sale_date BETWEEN $2 AND $3
ORDER BY s.sale_date DESC LIMIT $4
OFFSET $5;


-- name: GetRankedSalesManagers :many
-- get the ranked sales managers by their total sales divided by their sales goal amount for the given period.
WITH goal_sales AS (SELECT sm.sales_manager_id        AS sales_manager_id,
                           sm.first_name              AS first_name,
                           sm.last_name               AS last_name,
                           smg.amount                 AS sale_goal,
                           COALESCE(SUM(s.amount), 0) AS total_sales_sum
                    FROM sales_managers_view sm
                             LEFT JOIN
                         sales_manager_goals smg
                         ON sm.sales_manager_id = smg.sales_manager_id
                             AND smg.from_date = $1
                             AND smg.to_date = $2
                             LEFT JOIN
                         sales s
                         ON sm.sales_manager_id = s.sales_manager_id
                             AND s.sale_date BETWEEN $1 AND $2
                    GROUP BY sm.sales_manager_id,
                             smg.amount),
     rankings AS (SELECT *,
                         CASE
                             WHEN sale_goal = 0 THEN 0
                             ELSE total_sales_sum::decimal / sale_goal:: decimal
END
AS ratio,
        RANK() OVER (ORDER BY CASE
                             WHEN sale_goal = 0 THEN 0
                             ELSE total_sales_sum::decimal / sale_goal::decimal
                         END DESC) AS rating_position
    FROM
        goal_sales
)
SELECT sales_manager_id,
       sale_goal,
       total_sales_sum,
       ratio,
       rating_position
FROM rankings
ORDER BY rating_position LIMIT $3
OFFSET $4;

-- name: GetSalesCount :one
SELECT COUNT(*)
FROM sales
WHERE sales_manager_id = $1;

