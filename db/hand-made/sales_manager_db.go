package hand_made

import (
	"context"
	"database/sql"
)

type CustomQuerier interface {
	GetSalesManagerYearStatistic(ctx context.Context, arg GetSalesManagerYearStatisticParams) (*GetSalesManagerYearStatisticRow, error)
	GetBranchYearStatistic(ctx context.Context, arg GetBranchYearStatisticParams) ([]GetBranchYearStatisticRow, error)
	GetBranchRankedSalesManagers(ctx context.Context, arg GetBranchRankedSalesManagersParams) ([]GetRankedSalesManagersRow, error)
	GetBranchSumByType(ctx context.Context, arg GetBranchSumByTypeParams) (GetBranchSumByTypeRow, error)
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
SELECT SUM(amount) AS total_sales
FROM sales
WHERE EXTRACT(MONTH FROM sale_date) = $1
  AND EXTRACT(YEAR FROM sale_date) = $2
  AND sale_type_id = $3;
`

type GetSalesManagerYearStatisticParams struct {
	SalesManagerID int32 `json:"sales_manager_id"`
	Year           int32 `json:"year"`
	Month          int32 `json:"month"`
}

type GetSalesManagerYearStatisticRow struct {
	TotalAmount int64 `json:"total_amount"`
}

func (d DBCustomQuerier) GetSalesManagerYearStatistic(ctx context.Context, arg GetSalesManagerYearStatisticParams) (*GetSalesManagerYearStatisticRow, error) {
	rows, err := d.db.QueryContext(ctx, getSalesManagerYearStatistic, arg.SalesManagerID, arg.Year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items GetSalesManagerYearStatisticRow
	for rows.Next() {
		var i GetSalesManagerYearStatisticRow
		if err := rows.Scan(&i.TotalAmount); err != nil {
			return nil, err
		}
		return &items, nil
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &items, nil
}
