// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: sale_types.sql

package generated

import (
	"context"
)

const createSaleType = `-- name: CreateSaleType :one
INSERT INTO sale_types (title, description)
VALUES ($1, $2)
RETURNING id
`

type CreateSaleTypeParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (q *Queries) CreateSaleType(ctx context.Context, arg CreateSaleTypeParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, createSaleType, arg.Title, arg.Description)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const getSaleTypeById = `-- name: GetSaleTypeById :one
SELECT id, title, description, created_at, color, gravity, value_type FROM sale_types
WHERE id = $1
`

func (q *Queries) GetSaleTypeById(ctx context.Context, id int32) (SaleType, error) {
	row := q.db.QueryRowContext(ctx, getSaleTypeById, id)
	var i SaleType
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.CreatedAt,
		&i.Color,
		&i.Gravity,
		&i.ValueType,
	)
	return i, err
}

const getSalesTypes = `-- name: GetSalesTypes :many
SELECT id, title, description, created_at, color, gravity, value_type FROM sale_types
`

func (q *Queries) GetSalesTypes(ctx context.Context) ([]SaleType, error) {
	rows, err := q.db.QueryContext(ctx, getSalesTypes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SaleType
	for rows.Next() {
		var i SaleType
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.CreatedAt,
			&i.Color,
			&i.Gravity,
			&i.ValueType,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
