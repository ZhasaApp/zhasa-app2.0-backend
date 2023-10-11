package brand

import (
	"context"
	"database/sql"
	"fmt"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/statistic"
)

type GetBranchBrandGoalFunc func(branchBrand int32, saleTypeId int32, period statistic.Period) (int64, error)

func NewGetBranchBrandGoalFunc(ctx context.Context, store generated.BranchStore) GetBranchBrandGoalFunc {
	return func(branchBrand int32, saleTypeId int32, period statistic.Period) (int64, error) {
		from, to := period.ConvertToTime()

		goal, err := store.GetBranchBrandGoalByGivenDateRange(ctx, generated.GetBranchBrandGoalByGivenDateRangeParams{
			BranchBrand: branchBrand,
			FromDate:    from,
			ToDate:      to,
			SaleTypeID:  saleTypeId,
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
