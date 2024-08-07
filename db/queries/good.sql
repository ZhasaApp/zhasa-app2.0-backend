-- name: GetGoodsByBrandId :many
SELECT m.id, m.name, m.description
FROM goods m
         JOIN brand_goods bm ON m.id = bm.good_id
WHERE bm.brand_id = $1;

-- name: CreateGood :one
INSERT INTO goods (name, description)
VALUES ($1, $2) RETURNING id;

-- name: AddGoodToBrand :exec
INSERT INTO brand_goods (brand_id, good_id)
VALUES ($1, $2);

