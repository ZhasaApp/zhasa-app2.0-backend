// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: brands.sql

package generated

import (
	"context"
)

const addBrand = `-- name: AddBrand :exec
INSERT INTO brands (title, description)
VALUES ($1, $2)
`

type AddBrandParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (q *Queries) AddBrand(ctx context.Context, arg AddBrandParams) error {
	_, err := q.db.ExecContext(ctx, addBrand, arg.Title, arg.Description)
	return err
}

const getBranchBrand = `-- name: GetBranchBrand :one
SELECT bb.id AS branch_brand
FROM branch_brands bb
WHERE bb.branch_id = $1
  AND bb.brand_id = $2
`

type GetBranchBrandParams struct {
	BranchID int32 `json:"branch_id"`
	BrandID  int32 `json:"brand_id"`
}

func (q *Queries) GetBranchBrand(ctx context.Context, arg GetBranchBrandParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, getBranchBrand, arg.BranchID, arg.BrandID)
	var branch_brand int32
	err := row.Scan(&branch_brand)
	return branch_brand, err
}

const getBranchBrands = `-- name: GetBranchBrands :many
SELECT b.id, b.title
FROM branch_brands bb
         JOIN brands b ON bb.brand_id = b.id
WHERE bb.branch_id = $1
`

type GetBranchBrandsRow struct {
	ID    int32  `json:"id"`
	Title string `json:"title"`
}

func (q *Queries) GetBranchBrands(ctx context.Context, branchID int32) ([]GetBranchBrandsRow, error) {
	rows, err := q.db.QueryContext(ctx, getBranchBrands, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBranchBrandsRow
	for rows.Next() {
		var i GetBranchBrandsRow
		if err := rows.Scan(&i.ID, &i.Title); err != nil {
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

const getBrands = `-- name: GetBrands :many
SELECT b.id, b.title, b.description
FROM brands b LIMIT $1
OFFSET $2
`

type GetBrandsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetBrandsRow struct {
	ID          int32  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (q *Queries) GetBrands(ctx context.Context, arg GetBrandsParams) ([]GetBrandsRow, error) {
	rows, err := q.db.QueryContext(ctx, getBrands, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBrandsRow
	for rows.Next() {
		var i GetBrandsRow
		if err := rows.Scan(&i.ID, &i.Title, &i.Description); err != nil {
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

const getUserBrands = `-- name: GetUserBrands :many
SELECT b.id, b.title, b.description
FROM brands b
         JOIN user_brands ub ON b.id = ub.brand_id
WHERE ub.user_id = $1
`

type GetUserBrandsRow struct {
	ID          int32  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (q *Queries) GetUserBrands(ctx context.Context, userID int32) ([]GetUserBrandsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserBrands, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserBrandsRow
	for rows.Next() {
		var i GetUserBrandsRow
		if err := rows.Scan(&i.ID, &i.Title, &i.Description); err != nil {
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

const updateBrand = `-- name: UpdateBrand :exec
UPDATE brands
SET title = $1, description = $2
WHERE id = $3
`

type UpdateBrandParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ID          int32  `json:"id"`
}

func (q *Queries) UpdateBrand(ctx context.Context, arg UpdateBrandParams) error {
	_, err := q.db.ExecContext(ctx, updateBrand, arg.Title, arg.Description, arg.ID)
	return err
}
