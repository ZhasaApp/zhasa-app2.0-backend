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
FROM brands b
         JOIN user_brands ub ON b.id = ub.brand_id
WHERE ub.user_id = $1;

-- name: GetBranchBrand :one
SELECT bb.id AS branch_brand
FROM branch_brands bb
WHERE bb.branch_id = $1
  AND bb.brand_id = $2;

-- name: AddBrand :exec
INSERT INTO brands (title, description)
VALUES ($1, $2);

-- name: UpdateBrand :exec
UPDATE brands
SET title = $1, description = $2
WHERE id = $3;
