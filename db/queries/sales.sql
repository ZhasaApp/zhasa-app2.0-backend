-- name: AddSaleOrReplace :one
-- add sale into sales by given sale_type_id, amount, date, user_id and on conflict replace
INSERT INTO sales (user_id, sale_date, amount, sale_type_id, description)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: AddSaleToBrand :one
INSERT INTO sales_brands (sale_id, brand_id)
VALUES ($1, $2) RETURNING *;

-- name: GetSaleSumByManagerByTypeByBrand :one
-- get the sales sums for a specific sales manager and each sale type within the given period.
SELECT st.id         AS sale_type_id,
       st.title      AS sale_type_title,
       SUM(s.amount) AS total_sales
FROM sale_types st
         JOIN sales s ON st.id = s.sale_type_id
         JOIN sales_brands sb ON sb.sale_id = s.id
WHERE s.user_id = $1
  AND s.sale_date BETWEEN $2 AND $3
  AND st.id = $4
  AND sb.brand_id = $5
GROUP BY st.id, sb.brand_id
ORDER BY st.id ASC;

-- Assuming you also have a sales table as previously discussed.
-- name: GetSaleSumByBranchByTypeByBrand :many
-- Assuming you also have a sales table as previously discussed.
SELECT b.id          AS branch_id,
       b.title       AS branch_title,
       br.id         AS brand_id,
       br.title      AS brand_title,
       st.id         AS sale_type_id,
       st.title      AS sale_type_title,
       SUM(s.amount) AS total_sales
FROM sales s

-- Join with relevant tables
         JOIN users sm ON s.user_id = sm.id
         JOIN sale_types st ON s.sale_type_id = st.id
         JOIN branch_users_roles bur ON sm.user_id = bur.user_role_id
         JOIN branches b ON bur.branch_id = b.id
         JOIN branch_brands bb ON b.id = bb.branch_id
         JOIN brands br ON bb.brand_id = br.id

WHERE s.sale_date BETWEEN $1 AND $2
  AND b.id = $3
  AND br.id = $4

GROUP BY b.id, br.id, st.id
ORDER BY b.id, br.id, st.id;

-- name: GetSaleSumByUserIdBrandIdPeriodSaleTypeId :one
SELECT
    SUM(s.amount) AS total_sales
FROM
    sales s
        JOIN
    sales_brands sb ON s.id = sb.sale_id
        JOIN
    user_brands ub ON ub.brand_id = sb.brand_id
        JOIN
    users u ON u.id = ub.user_id
WHERE
        u.id = $1           -- user_id parameter
  AND
        sb.brand_id = $2     -- brand_id parameter
  AND
    s.sale_date BETWEEN $3 AND $4   -- from and to date parameters
  AND
        s.sale_type_id = $5   -- sale_type_id parameter
;

