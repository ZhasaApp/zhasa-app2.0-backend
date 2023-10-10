package hand_made

import (
	"context"
	"time"
)

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

func (q DBCustomQuerier) GetSaleSumByUserIdBrandIdPeriodSaleTypeId(ctx context.Context, arg GetSaleSumByUserIdBrandIdPeriodSaleTypeIdParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getSaleSumByUserIdBrandIdPeriodSaleTypeId,
		arg.ID,
		arg.BrandID,
		arg.SaleDate,
		arg.SaleDate_2,
		arg.SaleTypeID,
	)
	var totalSales int64
	err := row.Scan(&totalSales)
	return totalSales, err
}
