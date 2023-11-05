-- name: GetUserBrandGoal :one
SELECT COALESCE(goals.value, 0)
FROM user_brand_sale_type_goals goals
WHERE goals.user_id = $1
  AND goals.brand_id = $2
  AND goals.sale_type_id = $3
  AND goals.from_date = $4
  AND goals.to_date = $5;

-- name: GetUserBrand :one
SELECT ub.id AS user_brand
FROM user_brands ub
WHERE ub.user_id = $1
  AND ub.brand_id = $2;

-- name: InsertUserBrandRatio :exec
INSERT INTO user_brand_ratio (user_id, brand_id, ratio, from_date, to_date)
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (user_id, brand_id, from_date, to_date)
DO
UPDATE SET ratio = EXCLUDED.ratio;

-- name: GetUserRank :one
WITH RankedUsers AS (SELECT user_id,
                            brand_id,
                            ratio,
                            ROW_NUMBER() OVER (ORDER BY ratio DESC) as rank
                     FROM user_brand_ratio
                     WHERE brand_id = $1
                       AND from_date = $2
                       AND to_date = $3)

SELECT rank
FROM RankedUsers
WHERE user_id = $4;

-- name: DeleteUserBrandByUserId :exec
DELETE FROM user_brands
WHERE user_id = $1;
