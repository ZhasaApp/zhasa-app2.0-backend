package hand_made

import (
	"context"
	"database/sql"
)

type CustomQuerier interface {
	GetSalesManagerYearStatistic(ctx context.Context, arg GetSalesManagerYearStatisticParams) ([]GetSalesManagerYearStatisticRow, error)
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
SELECT 
    EXTRACT(MONTH FROM s.sale_date) AS month_number,
    s.sale_type_id as sale_type,
    COALESCE(g.amount, 0) AS goal,
    COALESCE(SUM(s.amount),0) AS total_amount
FROM 
    sales s
LEFT JOIN 
    sales_manager_goals_by_types g ON g.sales_manager_id = s.sales_manager_id AND g.type_id = s.sale_type_id
WHERE 
    s.sales_manager_id = $1 
    AND EXTRACT(YEAR FROM s.sale_date) = $2
GROUP BY 
    month_number,
    s.sale_type_id,
    goal;
`

type GetSalesManagerYearStatisticParams struct {
	SalesManagerID int32 `json:"sales_manager_id"`
	Year           int32 `json:"year"`
}

type GetSalesManagerYearStatisticRow struct {
	SaleType    int32 `json:"sale_type"`
	MonthNumber int32 `json:"month_number"`
	TotalAmount int64 `json:"total_amount"`
	Goal        int64 `json:"goal"`
}

func (d DBCustomQuerier) GetSalesManagerYearStatistic(ctx context.Context, arg GetSalesManagerYearStatisticParams) ([]GetSalesManagerYearStatisticRow, error) {
	rows, err := d.db.QueryContext(ctx, getSalesManagerYearStatistic, arg.SalesManagerID, arg.Year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSalesManagerYearStatisticRow
	for rows.Next() {
		var i GetSalesManagerYearStatisticRow
		if err := rows.Scan(&i.MonthNumber, &i.SaleType, &i.TotalAmount, &i.Goal); err != nil {
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
