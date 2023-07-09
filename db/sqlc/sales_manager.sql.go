// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: sales_manager.sql

package generated

import (
	"context"
	"time"
)

const addSaleOrReplace = `-- name: AddSaleOrReplace :one
INSERT INTO sales (sales_manager_id, sale_date, amount, sale_type_id, description)
VALUES ($1, $2, $3, $4, $5) RETURNING id, sales_manager_id, sale_date, amount, sale_type_id, description, created_at
`

type AddSaleOrReplaceParams struct {
	SalesManagerID int32     `json:"sales_manager_id"`
	SaleDate       time.Time `json:"sale_date"`
	Amount         int64     `json:"amount"`
	SaleTypeID     int32     `json:"sale_type_id"`
	Description    string    `json:"description"`
}

// add sale into sales by given sale_type_id, amount, date, sales_manager_id and on conflict replace
func (q *Queries) AddSaleOrReplace(ctx context.Context, arg AddSaleOrReplaceParams) (Sale, error) {
	row := q.db.QueryRowContext(ctx, addSaleOrReplace,
		arg.SalesManagerID,
		arg.SaleDate,
		arg.Amount,
		arg.SaleTypeID,
		arg.Description,
	)
	var i Sale
	err := row.Scan(
		&i.ID,
		&i.SalesManagerID,
		&i.SaleDate,
		&i.Amount,
		&i.SaleTypeID,
		&i.Description,
		&i.CreatedAt,
	)
	return i, err
}

const createSalesManager = `-- name: CreateSalesManager :exec
INSERT INTO sales_managers (user_id, branch_id)
VALUES ($1, $2)
`

type CreateSalesManagerParams struct {
	UserID   int32 `json:"user_id"`
	BranchID int32 `json:"branch_id"`
}

func (q *Queries) CreateSalesManager(ctx context.Context, arg CreateSalesManagerParams) error {
	_, err := q.db.ExecContext(ctx, createSalesManager, arg.UserID, arg.BranchID)
	return err
}

const getManagerSales = `-- name: GetManagerSales :many
SELECT id, sale_type_id, description, sale_date, amount
FROM sales s
WHERE s.sales_manager_id = $1
ORDER BY s.sale_date DESC LIMIT $2
OFFSET $3
`

type GetManagerSalesParams struct {
	SalesManagerID int32 `json:"sales_manager_id"`
	Limit          int32 `json:"limit"`
	Offset         int32 `json:"offset"`
}

type GetManagerSalesRow struct {
	ID          int32     `json:"id"`
	SaleTypeID  int32     `json:"sale_type_id"`
	Description string    `json:"description"`
	SaleDate    time.Time `json:"sale_date"`
	Amount      int64     `json:"amount"`
}

func (q *Queries) GetManagerSales(ctx context.Context, arg GetManagerSalesParams) ([]GetManagerSalesRow, error) {
	rows, err := q.db.QueryContext(ctx, getManagerSales, arg.SalesManagerID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetManagerSalesRow
	for rows.Next() {
		var i GetManagerSalesRow
		if err := rows.Scan(
			&i.ID,
			&i.SaleTypeID,
			&i.Description,
			&i.SaleDate,
			&i.Amount,
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

const getManagerSalesByPeriod = `-- name: GetManagerSalesByPeriod :many
SELECT id, sale_type_id, description, sale_date, amount
FROM sales s
WHERE s.sales_manager_id = $1
  AND s.sale_date BETWEEN $2 AND $3
ORDER BY s.sale_date DESC LIMIT $4
OFFSET $5
`

type GetManagerSalesByPeriodParams struct {
	SalesManagerID int32     `json:"sales_manager_id"`
	SaleDate       time.Time `json:"sale_date"`
	SaleDate_2     time.Time `json:"sale_date_2"`
	Limit          int32     `json:"limit"`
	Offset         int32     `json:"offset"`
}

type GetManagerSalesByPeriodRow struct {
	ID          int32     `json:"id"`
	SaleTypeID  int32     `json:"sale_type_id"`
	Description string    `json:"description"`
	SaleDate    time.Time `json:"sale_date"`
	Amount      int64     `json:"amount"`
}

func (q *Queries) GetManagerSalesByPeriod(ctx context.Context, arg GetManagerSalesByPeriodParams) ([]GetManagerSalesByPeriodRow, error) {
	rows, err := q.db.QueryContext(ctx, getManagerSalesByPeriod,
		arg.SalesManagerID,
		arg.SaleDate,
		arg.SaleDate_2,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetManagerSalesByPeriodRow
	for rows.Next() {
		var i GetManagerSalesByPeriodRow
		if err := rows.Scan(
			&i.ID,
			&i.SaleTypeID,
			&i.Description,
			&i.SaleDate,
			&i.Amount,
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

const getOrderedSalesManagers = `-- name: GetOrderedSalesManagers :many
SELECT
    v.sales_manager_id,
    v.first_name,
    v.last_name,
    v.avatar_url,
    v.branch_title,
    v.branch_id,
    v.user_id,
    COALESCE(r.ratio, 0.0) AS ratio
FROM
    sales_managers_view v
        LEFT JOIN
    sales_manager_goals_ratio_by_period r ON v.sales_manager_id = r.sales_manager_id
        AND r.from_date >= $1 AND r.to_date <= $2
ORDER BY
    ratio DESC
    LIMIT $3 OFFSET $4
`

type GetOrderedSalesManagersParams struct {
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
	Limit    int32     `json:"limit"`
	Offset   int32     `json:"offset"`
}

type GetOrderedSalesManagersRow struct {
	SalesManagerID int32   `json:"sales_manager_id"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	AvatarUrl      string  `json:"avatar_url"`
	BranchTitle    string  `json:"branch_title"`
	BranchID       int32   `json:"branch_id"`
	UserID         int32   `json:"user_id"`
	Ratio          float64 `json:"ratio"`
}

func (q *Queries) GetOrderedSalesManagers(ctx context.Context, arg GetOrderedSalesManagersParams) ([]GetOrderedSalesManagersRow, error) {
	rows, err := q.db.QueryContext(ctx, getOrderedSalesManagers,
		arg.FromDate,
		arg.ToDate,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOrderedSalesManagersRow
	for rows.Next() {
		var i GetOrderedSalesManagersRow
		if err := rows.Scan(
			&i.SalesManagerID,
			&i.FirstName,
			&i.LastName,
			&i.AvatarUrl,
			&i.BranchTitle,
			&i.BranchID,
			&i.UserID,
			&i.Ratio,
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

const getOrderedSalesManagersOfBranch = `-- name: GetOrderedSalesManagersOfBranch :many
SELECT
    v.sales_manager_id,
    v.first_name,
    v.last_name,
    v.avatar_url,
    v.branch_title,
    v.user_id,
    COALESCE(r.ratio, 0.0) AS ratio
FROM
    sales_managers_view v
        LEFT JOIN
    sales_manager_goals_ratio_by_period r ON v.sales_manager_id = r.sales_manager_id
        AND r.from_date >= $1 AND r.to_date <= $2
    AND v.branch_id = $3
ORDER BY
    ratio DESC
    LIMIT $4 OFFSET $5
`

type GetOrderedSalesManagersOfBranchParams struct {
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
	BranchID int32     `json:"branch_id"`
	Limit    int32     `json:"limit"`
	Offset   int32     `json:"offset"`
}

type GetOrderedSalesManagersOfBranchRow struct {
	SalesManagerID int32   `json:"sales_manager_id"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	AvatarUrl      string  `json:"avatar_url"`
	BranchTitle    string  `json:"branch_title"`
	UserID         int32   `json:"user_id"`
	Ratio          float64 `json:"ratio"`
}

func (q *Queries) GetOrderedSalesManagersOfBranch(ctx context.Context, arg GetOrderedSalesManagersOfBranchParams) ([]GetOrderedSalesManagersOfBranchRow, error) {
	rows, err := q.db.QueryContext(ctx, getOrderedSalesManagersOfBranch,
		arg.FromDate,
		arg.ToDate,
		arg.BranchID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOrderedSalesManagersOfBranchRow
	for rows.Next() {
		var i GetOrderedSalesManagersOfBranchRow
		if err := rows.Scan(
			&i.SalesManagerID,
			&i.FirstName,
			&i.LastName,
			&i.AvatarUrl,
			&i.BranchTitle,
			&i.UserID,
			&i.Ratio,
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

const getSMRatio = `-- name: GetSMRatio :one
SELECT ratio
FROM sales_manager_goals_ratio_by_period smgr
WHERE smgr.from_date = $1
  AND smgr.to_date = $2
  AND smgr.sales_manager_id = $3
`

type GetSMRatioParams struct {
	FromDate       time.Time `json:"from_date"`
	ToDate         time.Time `json:"to_date"`
	SalesManagerID int32     `json:"sales_manager_id"`
}

func (q *Queries) GetSMRatio(ctx context.Context, arg GetSMRatioParams) (float64, error) {
	row := q.db.QueryRowContext(ctx, getSMRatio, arg.FromDate, arg.ToDate, arg.SalesManagerID)
	var ratio float64
	err := row.Scan(&ratio)
	return ratio, err
}

const getSalesByDate = `-- name: GetSalesByDate :many
SELECT id, sales_manager_id, sale_date, amount, sale_type_id, description, created_at
from sales s
WHERE s.sale_date = $1
`

func (q *Queries) GetSalesByDate(ctx context.Context, saleDate time.Time) ([]Sale, error) {
	rows, err := q.db.QueryContext(ctx, getSalesByDate, saleDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Sale
	for rows.Next() {
		var i Sale
		if err := rows.Scan(
			&i.ID,
			&i.SalesManagerID,
			&i.SaleDate,
			&i.Amount,
			&i.SaleTypeID,
			&i.Description,
			&i.CreatedAt,
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

const getSalesCount = `-- name: GetSalesCount :one
SELECT COUNT(*)
FROM sales
WHERE sales_manager_id = $1
`

func (q *Queries) GetSalesCount(ctx context.Context, salesManagerID int32) (int64, error) {
	row := q.db.QueryRowContext(ctx, getSalesCount, salesManagerID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getSalesManagerByUserId = `-- name: GetSalesManagerByUserId :one
SELECT user_id, phone, first_name, last_name, avatar_url, sales_manager_id, branch_id, branch_title
from sales_managers_view s
WHERE s.user_id = $1
`

func (q *Queries) GetSalesManagerByUserId(ctx context.Context, userID int32) (SalesManagersView, error) {
	row := q.db.QueryRowContext(ctx, getSalesManagerByUserId, userID)
	var i SalesManagersView
	err := row.Scan(
		&i.UserID,
		&i.Phone,
		&i.FirstName,
		&i.LastName,
		&i.AvatarUrl,
		&i.SalesManagerID,
		&i.BranchID,
		&i.BranchTitle,
	)
	return i, err
}

const getSalesManagerGoalByGivenDateRangeAndSaleType = `-- name: GetSalesManagerGoalByGivenDateRangeAndSaleType :one
SELECT COALESCE(sg.amount, 0) AS goal_amount
FROM sales_manager_goals_by_types sg
WHERE sg.sales_manager_id = $1
  AND sg.from_date = $2
  AND sg.to_date = $3
  AND sg.type_id = $4
`

type GetSalesManagerGoalByGivenDateRangeAndSaleTypeParams struct {
	SalesManagerID int32     `json:"sales_manager_id"`
	FromDate       time.Time `json:"from_date"`
	ToDate         time.Time `json:"to_date"`
	TypeID         int32     `json:"type_id"`
}

func (q *Queries) GetSalesManagerGoalByGivenDateRangeAndSaleType(ctx context.Context, arg GetSalesManagerGoalByGivenDateRangeAndSaleTypeParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getSalesManagerGoalByGivenDateRangeAndSaleType,
		arg.SalesManagerID,
		arg.FromDate,
		arg.ToDate,
		arg.TypeID,
	)
	var goal_amount int64
	err := row.Scan(&goal_amount)
	return goal_amount, err
}

const getSalesManagerSumsByType = `-- name: GetSalesManagerSumsByType :one
SELECT st.id         AS sale_type_id,
       st.title      AS sale_type_title,
       SUM(s.amount) AS total_sales
FROM sale_types st
         JOIN sales s ON st.id = s.sale_type_id AND s.sales_manager_id = $1 AND s.sale_date BETWEEN $2 AND $3
WHERE st.id = $4
GROUP BY st.id
ORDER BY st.id ASC
`

type GetSalesManagerSumsByTypeParams struct {
	SalesManagerID int32     `json:"sales_manager_id"`
	SaleDate       time.Time `json:"sale_date"`
	SaleDate_2     time.Time `json:"sale_date_2"`
	ID             int32     `json:"id"`
}

type GetSalesManagerSumsByTypeRow struct {
	SaleTypeID    int32  `json:"sale_type_id"`
	SaleTypeTitle string `json:"sale_type_title"`
	TotalSales    int64  `json:"total_sales"`
}

// get the sales sums for a specific sales manager and each sale type within the given period.
func (q *Queries) GetSalesManagerSumsByType(ctx context.Context, arg GetSalesManagerSumsByTypeParams) (GetSalesManagerSumsByTypeRow, error) {
	row := q.db.QueryRowContext(ctx, getSalesManagerSumsByType,
		arg.SalesManagerID,
		arg.SaleDate,
		arg.SaleDate_2,
		arg.ID,
	)
	var i GetSalesManagerSumsByTypeRow
	err := row.Scan(&i.SaleTypeID, &i.SaleTypeTitle, &i.TotalSales)
	return i, err
}

const setSMRatio = `-- name: SetSMRatio :exec
INSERT INTO sales_manager_goals_ratio_by_period
    (from_date, to_date, ratio, sales_manager_id)
VALUES ($1, $2, $3, $4) ON CONFLICT (from_date, to_date, sales_manager_id)
DO
UPDATE SET ratio = EXCLUDED.ratio
`

type SetSMRatioParams struct {
	FromDate       time.Time `json:"from_date"`
	ToDate         time.Time `json:"to_date"`
	Ratio          float64   `json:"ratio"`
	SalesManagerID int32     `json:"sales_manager_id"`
}

func (q *Queries) SetSMRatio(ctx context.Context, arg SetSMRatioParams) error {
	_, err := q.db.ExecContext(ctx, setSMRatio,
		arg.FromDate,
		arg.ToDate,
		arg.Ratio,
		arg.SalesManagerID,
	)
	return err
}
