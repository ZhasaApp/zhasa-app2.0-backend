package brand

import (
	"context"
	"database/sql"
	"fmt"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/statistic"
	"zhasa2.0/user/entities"
)

type GetBranchBrandGoalFunc func(branchId, brandId, saleTypeId int32, period statistic.Period) (int64, error)

func NewGetBranchBrandGoalFunc(ctx context.Context, store generated.BranchStore) GetBranchBrandGoalFunc {
	return func(branchId, brandId, saleTypeId int32, period statistic.Period) (int64, error) {
		from, to := period.ConvertToTime()

		goal, err := store.GetBranchBrandGoalByGivenDateRange(ctx, generated.GetBranchBrandGoalByGivenDateRangeParams{
			BranchID:   branchId,
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

type GetBranchBrandGoalV2Func func(branchId, brandId, leadMeasureID int32, period statistic.Period) (int64, error)

func NewGetBranchBrandGoalV2Func(ctx context.Context, store generated.BranchStore) GetBranchBrandGoalV2Func {
	return func(branchId, brandId, leadMeasureID int32, period statistic.Period) (int64, error) {
		from, to := period.ConvertToTime()

		goal, err := store.GetBranchBrandGoalByGivenDateRangeV2(ctx, generated.GetBranchBrandGoalByGivenDateRangeV2Params{
			BranchID:      sql.NullInt32{Int32: branchId, Valid: true},
			BrandID:       sql.NullInt32{Int32: brandId, Valid: true},
			DateFrom:      from,
			DateTo:        to,
			LeadMeasureID: leadMeasureID,
			Type:          string(entities.GoalTypeBranchBrand),
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
