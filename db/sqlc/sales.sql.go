// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: sales.sql

package generated

import (
	"context"
	"time"
)

const addSaleOrReplace = `-- name: AddSaleOrReplace :one
INSERT INTO sales (user_id, sale_date, amount, sale_type_id, description)
VALUES ($1, $2, $3, $4, $5) RETURNING id, user_id, sale_date, amount, sale_type_id, description, created_at
`

type AddSaleOrReplaceParams struct {
	UserID      int32     `json:"user_id"`
	SaleDate    time.Time `json:"sale_date"`
	Amount      int64     `json:"amount"`
	SaleTypeID  int32     `json:"sale_type_id"`
	Description string    `json:"description"`
}

// add sale into sales by given sale_type_id, amount, date, user_id and on conflict replace
func (q *Queries) AddSaleOrReplace(ctx context.Context, arg AddSaleOrReplaceParams) (Sale, error) {
	row := q.db.QueryRowContext(ctx, addSaleOrReplace,
		arg.UserID,
		arg.SaleDate,
		arg.Amount,
		arg.SaleTypeID,
		arg.Description,
	)
	var i Sale
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SaleDate,
		&i.Amount,
		&i.SaleTypeID,
		&i.Description,
		&i.CreatedAt,
	)
	return i, err
}

const addSaleToBrand = `-- name: AddSaleToBrand :one
INSERT INTO sales_brands (sale_id, brand_id)
VALUES ($1, $2) RETURNING sale_id, brand_id
`

type AddSaleToBrandParams struct {
	SaleID  int32 `json:"sale_id"`
	BrandID int32 `json:"brand_id"`
}

func (q *Queries) AddSaleToBrand(ctx context.Context, arg AddSaleToBrandParams) (SalesBrand, error) {
	row := q.db.QueryRowContext(ctx, addSaleToBrand, arg.SaleID, arg.BrandID)
	var i SalesBrand
	err := row.Scan(&i.SaleID, &i.BrandID)
	return i, err
}

const deleteSale = `-- name: DeleteSale :exec
DELETE
FROM sales
WHERE id = $1
`

func (q *Queries) DeleteSale(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteSale, id)
	return err
}

const getSaleBrandBySaleId = `-- name: GetSaleBrandBySaleId :one
SELECT sb.brand_id, s.sale_date
FROM sales_brands sb
         JOIN sales s ON s.id = sb.sale_id
WHERE sb.sale_id = $1
`

type GetSaleBrandBySaleIdRow struct {
	BrandID  int32     `json:"brand_id"`
	SaleDate time.Time `json:"sale_date"`
}

func (q *Queries) GetSaleBrandBySaleId(ctx context.Context, saleID int32) (GetSaleBrandBySaleIdRow, error) {
	row := q.db.QueryRowContext(ctx, getSaleBrandBySaleId, saleID)
	var i GetSaleBrandBySaleIdRow
	err := row.Scan(&i.BrandID, &i.SaleDate)
	return i, err
}

const getSaleSumByBranchByTypeByBrand = `-- name: GetSaleSumByBranchByTypeByBrand :many
SELECT b.id          AS branch_id,
       b.title       AS branch_title,
       br.id         AS brand_id,
       br.title      AS brand_title,
       st.id         AS sale_type_id,
       st.title      AS sale_type_title,
       SUM(s.amount) AS total_sales
FROM sales s

         JOIN users sm ON s.user_id = sm.id
         JOIN sale_types st ON s.sale_type_id = st.id
         JOIN branch_users bu ON sm.user_id = bu.user_id
         JOIN branches b ON bur.branch_id = b.id
         JOIN branch_brands bb ON b.id = bb.branch_id
         JOIN brands br ON bb.brand_id = br.id

WHERE s.sale_date BETWEEN $1 AND $2
  AND b.id = $3
  AND br.id = $4

GROUP BY b.id, br.id, st.id
ORDER BY b.id, br.id, st.id
`

type GetSaleSumByBranchByTypeByBrandParams struct {
	SaleDate   time.Time `json:"sale_date"`
	SaleDate_2 time.Time `json:"sale_date_2"`
	ID         int32     `json:"id"`
	ID_2       int32     `json:"id_2"`
}

type GetSaleSumByBranchByTypeByBrandRow struct {
	BranchID      int32  `json:"branch_id"`
	BranchTitle   string `json:"branch_title"`
	BrandID       int32  `json:"brand_id"`
	BrandTitle    string `json:"brand_title"`
	SaleTypeID    int32  `json:"sale_type_id"`
	SaleTypeTitle string `json:"sale_type_title"`
	TotalSales    int64  `json:"total_sales"`
}

// Assuming you also have a sales table as previously discussed.
// Assuming you also have a sales table as previously discussed.
// Join with relevant tables
func (q *Queries) GetSaleSumByBranchByTypeByBrand(ctx context.Context, arg GetSaleSumByBranchByTypeByBrandParams) ([]GetSaleSumByBranchByTypeByBrandRow, error) {
	rows, err := q.db.QueryContext(ctx, getSaleSumByBranchByTypeByBrand,
		arg.SaleDate,
		arg.SaleDate_2,
		arg.ID,
		arg.ID_2,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSaleSumByBranchByTypeByBrandRow
	for rows.Next() {
		var i GetSaleSumByBranchByTypeByBrandRow
		if err := rows.Scan(
			&i.BranchID,
			&i.BranchTitle,
			&i.BrandID,
			&i.BrandTitle,
			&i.SaleTypeID,
			&i.SaleTypeTitle,
			&i.TotalSales,
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

const getSaleSumByUserIdBrandIdPeriodSaleTypeId = `-- name: GetSaleSumByUserIdBrandIdPeriodSaleTypeId :one
SELECT COALESCE(SUM(s.amount), 0) AS total_sales
FROM sales s
         JOIN
     sales_brands sb ON s.id = sb.sale_id
         JOIN
     user_brands ub ON ub.brand_id = sb.brand_id
         JOIN
     users u ON u.id = ub.user_id
WHERE u.id = $1                     -- user_id parameter
  AND sb.brand_id = $2              -- brand_id parameter
  AND s.sale_date BETWEEN $3 AND $4 -- from and to date parameters
  AND s.sale_type_id = $5 -- sale_type_id parameter
`

type GetSaleSumByUserIdBrandIdPeriodSaleTypeIdParams struct {
	ID         int32     `json:"id"`
	BrandID    int32     `json:"brand_id"`
	SaleDate   time.Time `json:"sale_date"`
	SaleDate_2 time.Time `json:"sale_date_2"`
	SaleTypeID int32     `json:"sale_type_id"`
}

func (q *Queries) GetSaleSumByUserIdBrandIdPeriodSaleTypeId(ctx context.Context, arg GetSaleSumByUserIdBrandIdPeriodSaleTypeIdParams) (interface{}, error) {
	row := q.db.QueryRowContext(ctx, getSaleSumByUserIdBrandIdPeriodSaleTypeId,
		arg.ID,
		arg.BrandID,
		arg.SaleDate,
		arg.SaleDate_2,
		arg.SaleTypeID,
	)
	var total_sales interface{}
	err := row.Scan(&total_sales)
	return total_sales, err
}

const getSalesByBrandId = `-- name: GetSalesByBrandId :many
SELECT s.id,
       s.user_id,
       s.sale_date,
       s.amount,
       s.sale_type_id,
       s.description
FROM sales s
         JOIN
     sales_brands sb ON s.id = sb.sale_id
WHERE sb.brand_id = $1
`

type GetSalesByBrandIdRow struct {
	ID          int32     `json:"id"`
	UserID      int32     `json:"user_id"`
	SaleDate    time.Time `json:"sale_date"`
	Amount      int64     `json:"amount"`
	SaleTypeID  int32     `json:"sale_type_id"`
	Description string    `json:"description"`
}

func (q *Queries) GetSalesByBrandId(ctx context.Context, brandID int32) ([]GetSalesByBrandIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getSalesByBrandId, brandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSalesByBrandIdRow
	for rows.Next() {
		var i GetSalesByBrandIdRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.SaleDate,
			&i.Amount,
			&i.SaleTypeID,
			&i.Description,
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

const getSalesByBrandIdAndUserId = `-- name: GetSalesByBrandIdAndUserId :many
SELECT s.id,
       s.user_id,
       s.sale_date,
       s.amount,
       s.sale_type_id,
       s.description,
       st.title AS sale_type_title,
       st.gravity,
       st.color,
       st.value_type
FROM sales s
         JOIN sales_brands sb ON s.id = sb.sale_id
         JOIN sale_types st ON s.sale_type_id = st.id
WHERE sb.brand_id = $1
  AND s.user_id = $2
ORDER BY s.sale_date DESC LIMIT $3
OFFSET $4
`

type GetSalesByBrandIdAndUserIdParams struct {
	BrandID int32 `json:"brand_id"`
	UserID  int32 `json:"user_id"`
	Limit   int32 `json:"limit"`
	Offset  int32 `json:"offset"`
}

type GetSalesByBrandIdAndUserIdRow struct {
	ID            int32     `json:"id"`
	UserID        int32     `json:"user_id"`
	SaleDate      time.Time `json:"sale_date"`
	Amount        int64     `json:"amount"`
	SaleTypeID    int32     `json:"sale_type_id"`
	Description   string    `json:"description"`
	SaleTypeTitle string    `json:"sale_type_title"`
	Gravity       int32     `json:"gravity"`
	Color         string    `json:"color"`
	ValueType     ValueType `json:"value_type"`
}

func (q *Queries) GetSalesByBrandIdAndUserId(ctx context.Context, arg GetSalesByBrandIdAndUserIdParams) ([]GetSalesByBrandIdAndUserIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getSalesByBrandIdAndUserId,
		arg.BrandID,
		arg.UserID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSalesByBrandIdAndUserIdRow
	for rows.Next() {
		var i GetSalesByBrandIdAndUserIdRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.SaleDate,
			&i.Amount,
			&i.SaleTypeID,
			&i.Description,
			&i.SaleTypeTitle,
			&i.Gravity,
			&i.Color,
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
