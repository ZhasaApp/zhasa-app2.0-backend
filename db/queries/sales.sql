-- name: AddSaleOrReplace :one
-- add sale into sales by given sale_type_id, amount, date, user_id and on conflict replace
INSERT INTO sales (user_id, sale_date, amount, sale_type_id, description)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: EditSale :one
UPDATE sales
SET user_id      = $1,
    sale_date    = $2,
    amount       = $3,
    sale_type_id = $4,
    description  = $5
WHERE id = $6 RETURNING *;


-- name: AddSaleToBrand :one
INSERT INTO sales_brands (sale_id, brand_id)
VALUES ($1, $2) ON CONFLICT (sale_id, brand_id)
DO
UPDATE
    SET sale_id = EXCLUDED.sale_id, brand_id = EXCLUDED.brand_id
    RETURNING *;

-- name: AddGoodToSale :exec
INSERT INTO sales_goods (sale_id, good_id)
VALUES ($1, $2) ON CONFLICT (sale_id)
DO
UPDATE
    SET good_id = EXCLUDED.good_id;

-- name: DeleteSale :exec
DELETE
FROM sales
WHERE id = $1;

-- Assuming you also have a sales table as previously discussed.
-- name: GetSaleSumByBranchByTypeByBrand :one
-- Assuming you also have a sales table as previously discussed.
SELECT b.id     AS branch_id,
       b.title  AS branch_title,
       br.id    AS brand_id,
       br.title AS brand_title,
       st.id    AS sale_type_id,
       st.title AS sale_type_title,
       COALESCE(SUM(s.amount), 0) ::bigint AS total_sales
FROM sales s

-- Join with relevant tables
         JOIN users sm ON s.user_id = sm.id
         JOIN sale_types st ON s.sale_type_id = st.id
         JOIN branch_users bu ON sm.user_id = bu.user_id
         JOIN branches b ON bu.branch_id = b.id
         JOIN branch_brands bb ON b.id = bb.branch_id
         JOIN brands br ON bb.brand_id = br.id

WHERE s.sale_date BETWEEN $1 AND $2
  AND b.id = $3
  AND br.id = $4
  AND st.id = $5
GROUP BY b.id, br.id, st.id
ORDER BY b.id, br.id, st.id;

-- name: GetSalesByBrandId :many
SELECT s.id,
       s.user_id,
       s.sale_date,
       s.amount,
       s.sale_type_id,
       s.description
FROM sales s
         JOIN
     sales_brands sb ON s.id = sb.sale_id
WHERE sb.brand_id = $1;

-- name: GetSalesByBrandIdAndUserId :many
SELECT s.id,
       s.user_id,
       s.sale_date,
       s.amount,
       s.sale_type_id,
       s.description,
       st.title AS sale_type_title,
       st.gravity,
       st.color,
       st.value_type
FROM sales s
         JOIN sales_brands sb ON s.id = sb.sale_id
         JOIN sale_types st ON s.sale_type_id = st.id
WHERE sb.brand_id = $1
  AND s.user_id = $2
ORDER BY s.sale_date DESC LIMIT $3
OFFSET $4;

-- name: GetSaleBrandBySaleId :one
SELECT sb.brand_id, s.sale_date
FROM sales_brands sb
         JOIN sales s ON s.id = sb.sale_id
WHERE sb.sale_id = $1;

-- name: GetSumByUserIdBrandIdPeriodSaleTypeId :one
SELECT COALESCE(SUM(s.amount), 0) ::bigint AS total_sales
FROM sales s
         JOIN
     sales_brands sb ON s.id = sb.sale_id AND sb.brand_id = $2 -- brand_id parameter
         JOIN
     user_brands ub ON ub.brand_id = sb.brand_id AND ub.user_id = s.user_id
WHERE s.user_id = $1      -- user_id parameter
  AND s.sale_type_id = $3 -- sale_type_id parameter
  AND s.sale_date BETWEEN $4 AND $5; -- from_date and to_date parameters
