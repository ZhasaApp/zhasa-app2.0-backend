package repository

import (
	"context"
	"database/sql"
	"fmt"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/statistic"
)

type GetBrandOverallGoalFunc func(brandId, saleTypeId int32, period statistic.Period) (int64, error)

func NewGetBrandOverallGoalFunc(ctx context.Context, store generated.BranchStore) GetBrandOverallGoalFunc {
	return func(brandId, saleTypeId int32, period statistic.Period) (int64, error) {
		from, to := period.ConvertToTime()

		goal, err := store.GetBrandOverallGoalByGivenDateRange(ctx, generated.GetBrandOverallGoalByGivenDateRangeParams{
			BrandID:    brandId,
			FromDate:   from,
			ToDate:     to,
			SaleTypeID: saleTypeId,
		})
		if err == sql.ErrNoRows {
			return 0, nil
		}
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		return goal, nil
	}
}
