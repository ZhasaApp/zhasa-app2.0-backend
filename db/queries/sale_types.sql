-- name: GetSalesTypes :many
SELECT * FROM sale_types;

-- name: GetSaleTypeById :one
SELECT * FROM sale_types
WHERE id = $1;

-- name: CreateSaleType :one
INSERT INTO sale_types (title, description)
VALUES ($1, $2)
RETURNING id;