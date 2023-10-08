-- name: GetUserBrandGoal :one
SELECT COALESCE(goals.value, 0)
FROM user_brand_sale_type_goals goals
WHERE goals.user_brand = $1
  AND goals.sale_type_id = $2
  AND goals.from_date = $3
  AND goals.from_date = $4;

-- name: GetUserBrand :one
SELECT ub.id AS user_brand
FROM user_brands ub
WHERE ub.user_id = $1
  AND ub.brand_id = $2;

-- name: InsertUserBrandRatio :exec
INSERT INTO user_brand_sale_type_ratio (user_id, brand_id, sale_type_id, ratio, from_date, to_date)
VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (user_id, brand_id, sale_type_id, from_date, to_date)
DO
UPDATE SET ratio = EXCLUDED.ratio;

-- name: GetUsersOrderedByRatioForGivenBrand :many
-- SELECT users for given brand ordered by ratio
SELECT u.id, u.first_name, u.last_name, r.ratio
FROM users u
         JOIN user_brand_sale_type_ratio r ON u.id = r.user_id
WHERE r.brand_id = $1
  AND r.from_date = $2
  AND r.to_date = $3
ORDER BY r.ratio DESC
OFFSET $4 LIMIT $5;

-- name: GetUserRank :one
WITH RankedUsers AS (SELECT user_id,
                            brand_id,
                            ratio,
                            ROW_NUMBER() OVER (ORDER BY ratio DESC) as rank
                     FROM user_brand_sale_type_ratio
                     WHERE brand_id = $1
                       AND from_date = $2
                       AND to_date = $3)

SELECT rank
FROM RankedUsers
WHERE user_id = $4;
