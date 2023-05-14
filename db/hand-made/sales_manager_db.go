package hand_made

import (
	"context"
	"database/sql"
)

type CustomQuerier interface {
	GetSalesManagerYearStatistic(ctx context.Context, arg GetSalesManagerYearStatisticParams) ([]GetSalesManagerYearStatisticRow, error)
	GetBranchYearStatistic(ctx context.Context, arg GetBranchYearStatisticParams) ([]GetBranchYearStatisticRow, error)
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
SELECT st.id AS sale_type,
       CAST(EXTRACT(MONTH FROM s.sale_date) AS INTEGER) AS month_number,
       SUM(s.amount) AS total_amount
FROM sales AS s
         JOIN sale_types AS st ON s.sale_type_id = st.id
WHERE s.sales_manager_id = $1
  AND DATE_PART('year', s.sale_date)::integer = $2
GROUP BY st.id, EXTRACT (MONTH FROM s.sale_date)
ORDER BY month_number, st.id
`

type GetSalesManagerYearStatisticParams struct {
	SalesManagerID int32 `json:"sales_manager_id"`
	Year           int32 `json:"year"`
}

type GetSalesManagerYearStatisticRow struct {
	SaleType    int32 `json:"sale_type"`
	MonthNumber int32 `json:"month_number"`
	TotalAmount int64 `json:"total_amount"`
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
		if err := rows.Scan(&i.SaleType, &i.MonthNumber, &i.TotalAmount); err != nil {
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
