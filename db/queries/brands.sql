
-- name: GetBranchBrands :many
SELECT b.id, b.title
FROM branch_brands bb
         JOIN brands b ON bb.brand_id = b.id
WHERE bb.branch_id = $1;