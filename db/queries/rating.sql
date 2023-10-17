-- name: GetUsersOrderedByRatioForGivenBrand :many
-- SELECT distinct users for given brand ordered by ratio and limited by offset and limit and if there is no any user with ratio let ratio be 0
SELECT u.id,
       u.first_name,
       u.last_name,
       u.avatar_url,
       COALESCE(r.ratio, 0) AS ratio,
       b.title              AS branch_title,
       b.id                 AS branch_id
FROM user_avatar_view u
         JOIN
     branch_users bu ON u.id = bu.user_id
         JOIN
     branches b ON bu.branch_id = b.id
         JOIN
     user_brands ub ON u.id = ub.user_id AND ub.brand_id = $1
         LEFT JOIN
     user_brand_ratio r ON u.id = r.user_id AND r.from_date = $2 AND r.to_date = $3
         JOIN user_roles ur ON u.id = ur.user_id AND ur.role_id = $6
WHERE (r.brand_id = $1 OR r.brand_id IS NULL)
ORDER BY r.ratio
OFFSET $4 LIMIT $5;

-- name: GetBranchUsersOrderedByRatioForGivenBrand :many
-- SELECT distinct users for given brand ordered by ratio and limited by offset and limit and if there is no any user with ratio let ratio be 0
SELECT u.id,
       u.first_name,
       u.last_name,
       u.avatar_url,
       COALESCE(r.ratio, 0) AS ratio,
       b.title              AS branch_title,
       b.id                 AS branch_id
FROM user_avatar_view u
         JOIN
     branch_users bu ON u.id = bu.user_id
         JOIN
     branches b ON bu.branch_id = b.id
         JOIN
     user_brands ub ON u.id = ub.user_id AND ub.brand_id = $1
         LEFT JOIN
     user_brand_ratio r ON u.id = r.user_id AND r.from_date = $2 AND r.to_date = $3
         JOIN user_roles ur ON u.id = ur.user_id AND ur.role_id = $7
WHERE (r.brand_id = $1 OR r.brand_id IS NULL)
  AND b.id = $6
ORDER BY r.ratio
OFFSET $4 LIMIT $5;



