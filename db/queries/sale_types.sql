-- name: GetSalesTypes :many
SELECT * FROM sale_types;

-- name: GetSaleTypeById :one
SELECT * FROM sale_types
WHERE id = $1;