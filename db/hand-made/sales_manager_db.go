package hand_made

import (
	"context"
	"database/sql"
)

type CustomQuerier interface {
	GetSalesManagerYearStatistic(ctx context.Context, arg GetSalesManagerYearStatisticParams) (*GetSalesManagerYearStatisticRow, error)
	GetBranchYearStatistic(ctx context.Context, arg GetBranchYearStatisticParams) (*GetBranchYearStatisticRow, error)
	GetBranchSumByType(ctx context.Context, arg GetBranchSumByTypeParams) (GetBranchSumByTypeRow, error)
	GetPostsAndPostAuthors(ctx context.Context, arg GetPostsAndPostAuthorsParams) ([]GetPostsAndPostAuthorsRow, error)
}

func NewCustomQuerier(db *sql.DB) CustomQuerier {
	return DBCustomQuerier{
		db: db,
	}
}

type DBCustomQuerier struct {
	db *sql.DB
}

const getSalesManagerYearStatistic = `-- name: GetSalesManagerYearStatistic :many
SELECT COALESCE(SUM(amount),0) AS total_sales
FROM sales
WHERE EXTRACT(MONTH FROM sale_date) = $1
  AND EXTRACT(YEAR FROM sale_date) = $2
  AND sales_manager_id = $3
  AND sale_type_id = $4;
`

type GetSalesManagerYearStatisticParams struct {
	SalesManagerID int32 `json:"sales_manager_id"`
	TypeId         int32 `json:"type_id"`
	Year           int32 `json:"year"`
	Month          int32 `json:"month"`
}

type GetSalesManagerYearStatisticRow struct {
	TotalAmount int64 `json:"total_amount"`
}

func (d DBCustomQuerier) GetSalesManagerYearStatistic(ctx context.Context, arg GetSalesManagerYearStatisticParams) (*GetSalesManagerYearStatisticRow, error) {
	rows := d.db.QueryRowContext(ctx, getSalesManagerYearStatistic, arg.Month, arg.Year, arg.SalesManagerID, arg.TypeId)
	var i GetSalesManagerYearStatisticRow
	if err := rows.Scan(&i.TotalAmount); err != nil {
		return nil, err
	}
	return &i, nil
}
