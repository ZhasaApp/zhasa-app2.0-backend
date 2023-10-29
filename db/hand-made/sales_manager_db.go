package hand_made

import (
	"context"
)

const getSalesManagerYearStatistic = `-- name: GetSalesManagerYearStatistic :many
SELECT COALESCE(SUM(amount),0) AS total_sales
FROM sales JOIN sales_brands sb ON sales.id = sb.sale_id
WHERE EXTRACT(MONTH FROM sale_date) = $1
  AND EXTRACT(YEAR FROM sale_date) = $2
  AND user_id = $3
  AND sale_type_id = $4
AND sb.brand_id = $5;
`

type GetSalesManagerYearStatisticParams struct {
	SalesManagerID int32 `json:"sales_manager_id"`
	TypeId         int32 `json:"type_id"`
	Year           int32 `json:"year"`
	Month          int32 `json:"month"`
	BrandId        int32 `json:"brand_id"`
}

type GetSalesManagerYearStatisticRow struct {
	TotalAmount int64 `json:"total_amount"`
}

func (d DBCustomQuerier) GetSalesManagerYearStatistic(ctx context.Context, arg GetSalesManagerYearStatisticParams) (*GetSalesManagerYearStatisticRow, error) {
	rows := d.db.QueryRowContext(ctx, getSalesManagerYearStatistic, arg.Month, arg.Year, arg.SalesManagerID, arg.TypeId, arg.BrandId)
	var i GetSalesManagerYearStatisticRow
	if err := rows.Scan(&i.TotalAmount); err != nil {
		return nil, err
	}
	return &i, nil
}
