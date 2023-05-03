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
WITH sales_by_manager_type AS (SELECT sm.id         AS sales_manager_id,
                                      st.id         AS sale_type_id,
                                      SUM(s.amount) AS total_sales
                               FROM sales s
                                        INNER JOIN sales_managers sm ON s.sales_manager_id = sm.id
                                        INNER JOIN sale_types st ON s.sale_type_id = st.id
                               WHERE s.sale_date BETWEEN $1 AND $2
                                 AND sm.id = $3
                               GROUP BY sm.id,
                                        st.id)
SELECT smt.sales_manager_id,
       smt.sale_type_id,
       COALESCE(smt.total_sales, 0) AS total_sales
FROM sales_by_manager_type smt
ORDER BY smt.sale_type_id ASC;

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