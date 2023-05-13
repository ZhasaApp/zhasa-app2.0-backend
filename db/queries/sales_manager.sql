-- name: CreateSalesManager :exec
INSERT INTO sales_managers (user_id, branch_id)
VALUES ($1, $2);

-- name: GetRankedSalesManagers :many
-- get the ranked sales managers by their total sales divided by their sales goal amount for the given period.
WITH sales_summary AS (SELECT sm.id         AS sales_manager_id,
                              SUM(s.amount) AS total_sales_amount,
                              u.first_name  AS first_name,
                              u.last_name   AS last_name,
                              u.avatar_url  AS avatar_url
                       FROM sales s
                                INNER JOIN sales_managers sm ON s.sales_manager_id = sm.id
                                INNER JOIN user_avatar_view u ON sm.user_id = u.id
                       WHERE s.sale_date BETWEEN $1 AND $2
                       GROUP BY sm.id),
     goal_summary AS (SELECT sm.id     AS sales_manager_id,
                             sg.from_date,
                             sg.to_date,
                             sg.amount AS goal_amount
                      FROM sales_manager_goals sg
                               INNER JOIN sales_managers sm ON s.sales_manager_id = sm.id
                      WHERE sg.from_date = $1
                        AND sg.to_date = $2)
SELECT ss.sales_manager_id,
       ss.first_name,
       ss.last_name,
       ss.avatar_url,
       COALESCE(ss.total_sales_amount / NULLIF(smg.goal_amount, 0), 0) ::float AS ratio
FROM sales_summary ss
         LEFT JOIN goal_summary smg ON ss.sales_manager_id = smg.sales_manager_id
ORDER BY ratio DESC LIMIT $3
OFFSET $4;

-- name: GetSalesManagerSumsByType :many
-- get the sales sums for a specific sales manager and each sale type within the given period.
SELECT st.id         AS sale_type_id,
       st.title      AS sale_type_title,
       SUM(s.amount) AS total_sales
FROM sale_types st
         JOIN sales s ON st.id = s.sale_type_id AND s.sales_manager_id = $1 AND s.sale_date BETWEEN $2 AND $3
GROUP BY st.id
ORDER BY st.id ASC;



-- name: AddSaleOrReplace :exec
-- add sale into sales by given sale_type_id, amount, date, sales_manager_id and on conflict replace
INSERT INTO sales (sales_manager_id, sale_date, amount, sale_type_id, description)
VALUES ($1, $2, $3, $4, $5);

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

-- name: GetSalesManagerYearStatistic :many
SELECT st.id AS sale_type,
       CAST(EXTRACT(MONTH FROM s.sale_date) AS INTEGER) AS month_number,
       SUM(s.amount) AS total_amount
FROM sales AS s
         JOIN sale_types AS st ON s.sale_type_id = st.id
WHERE s.sales_manager_id = $1
  AND DATE_PART('year', s.sale_date)::integer = $2
GROUP BY st.id, EXTRACT (MONTH FROM s.sale_date)
ORDER BY month_number, st.id;
