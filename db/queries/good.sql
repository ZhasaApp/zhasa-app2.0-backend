-- name: GetGoodsByBrandId :many
SELECT m.id, m.name, m.description
FROM goods m
         JOIN brand_goods bm ON m.id = bm.good_id
WHERE bm.brand_id = $1
  AND NOT EXISTS (
    SELECT 1
    FROM disabled_goods dg
    WHERE dg.good_id = m.id
);

-- name: CreateGood :one
INSERT INTO goods (name, description)
VALUES ($1, $2) RETURNING id;

-- name: AddGoodToBrand :exec
INSERT INTO brand_goods (brand_id, good_id)
VALUES ($1, $2);

-- name: GetGoodBySaleId :one
SELECT g.id, g.name, g.description
FROM goods g
         JOIN sales_goods sg ON g.id = sg.good_id
WHERE sg.sale_id = $1;

-- name: DisableGood :exec
INSERT INTO disabled_goods (good_id)
VALUES ($1);
