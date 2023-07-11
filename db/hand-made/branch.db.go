package hand_made

import (
	"context"
	"time"
)

const getBranchYearStatistic = `-- name: GetBranchYearStatistic :many
SELECT COALESCE(SUM(amount), 0) AS total_sales
FROM sales
JOIN sales_managers ON sales.sales_manager_id = sales_managers.id
JOIN branches ON sales_managers.branch_id = branches.id
WHERE EXTRACT(MONTH FROM sale_date) = $1
  AND EXTRACT(YEAR FROM sale_date) = $2
  AND branches.id = $3
  AND sale_type_id = $4;
`

type GetBranchYearStatisticParams struct {
	BranchId int32 `json:"branch_id"`
	TypeId   int32 `json:"type_id"`
	Year     int32 `json:"year"`
	Month    int32 `json:"month"`
}
type GetBranchYearStatisticRow struct {
	TotalAmount int64 `json:"total_amount"`
}

func (d DBCustomQuerier) GetBranchYearStatistic(ctx context.Context, arg GetBranchYearStatisticParams) (*GetBranchYearStatisticRow, error) {
	rows := d.db.QueryRowContext(ctx, getBranchYearStatistic, arg.Year, arg.Month, arg.BranchId, arg.TypeId)

	var i GetBranchYearStatisticRow

	if err := rows.Scan(&i.TotalAmount); err != nil {
		return nil, err
	}
	return &i, nil
}

const getBranchSumByType = `-- name: GetBranchSumByType :one
SELECT COALESCE(SUM(s.amount), 0) AS total_sales
FROM sales AS s
         JOIN sales_managers AS sm ON s.sales_manager_id = sm.id
         JOIN branches AS b ON sm.branch_id = b.id
         JOIN sale_types AS st ON s.sale_type_id = st.id
WHERE b.id = $1
  AND st.id = $4
  AND s.sale_date BETWEEN $2 AND $3
`

type GetBranchSumByTypeParams struct {
	BranchID   int32     `json:"branch_id"`
	SaleDate   time.Time `json:"sale_date"`
	SaleDate_2 time.Time `json:"sale_date_2"`
	ID         int32     `json:"id"`
}

type GetBranchSumByTypeRow struct {
	TotalSales int64 `json:"total_sales"`
}

func (d DBCustomQuerier) GetBranchSumByType(ctx context.Context, arg GetBranchSumByTypeParams) (GetBranchSumByTypeRow, error) {
	row := d.db.QueryRowContext(ctx, getBranchSumByType,
		arg.BranchID,
		arg.SaleDate,
		arg.SaleDate_2,
		arg.ID,
	)
	var i GetBranchSumByTypeRow
	err := row.Scan(
		&i.TotalSales,
	)
	return i, err
}
