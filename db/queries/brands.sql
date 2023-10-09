-- name: GetBranchBrands :many
SELECT b.id, b.title
FROM branch_brands bb
         JOIN brands b ON bb.brand_id = b.id
WHERE bb.branch_id = $1;

-- name: GetBrands :many
SELECT b.id, b.title, b.description
FROM brands b LIMIT $1
OFFSET $2;

-- name: GetUserBrands :many
SELECT b.id, b.title, b.description
FROM brands b JOIN user_brands ub ON b.id = ub.brand_id
WHERE ub.user_id = $1;
