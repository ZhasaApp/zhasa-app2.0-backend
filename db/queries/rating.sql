-- name: GetUsersOrderedByRatioForGivenBrand :many
-- SELECT distinct users for given brand ordered by ratio and limited by offset and limit and if there is no any user with ratio let ratio be 0
SELECT u.id,
       u.first_name,
       u.last_name,
       u.avatar_url,
       COALESCE(r.ratio, 0) AS ratio,
       b.title              AS branch_title,
       b.id                 as branch_id
FROM user_avatar_view u
         LEFT JOIN user_brand_ratio r
         JOIN branch_users bu ON u.id = bu.user_id
         JOIN branches b ON bu.branch_id = b.id
              ON u.id = r.user_id
WHERE r.brand_id = $1
  AND r.from_date = $2
  AND r.to_date = $3
ORDER BY r.ratio DESC
OFFSET $4 LIMIT $5;
