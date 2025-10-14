package repository

import (
	"context"
	"database/sql"
	"fmt"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/statistic"
	"zhasa2.0/user/entities"
)

type SetBranchBrandSaleTypeGoal func(branchId, brandId, saleTypeId int32, goal int64, period statistic.Period) error

func NewSetBranchGoalFunc(ctx context.Context, store generated.BranchStore) SetBranchBrandSaleTypeGoal {
	return func(branchId, brandId, saleTypeId int32, goal int64, period statistic.Period) error {
		from, to := period.ConvertToTime()
		err := store.SetBranchBrandGoal(ctx, generated.SetBranchBrandGoalParams{
			BranchID:   branchId,
			BrandID:    brandId,
			SaleTypeID: saleTypeId,
			Value:      goal,
			FromDate:   from,
			ToDate:     to,
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	}
}

type SetBranchBrandSaleTypeGoalV2 func(branchId, brandId, leadMeasureID int32, goal int64, period statistic.Period) error

func NewSetBranchGoalV2Func(ctx context.Context, store generated.BranchStore) SetBranchBrandSaleTypeGoalV2 {
	return func(branchId, brandId, leadMeasureID int32, goal int64, period statistic.Period) error {
		from, to := period.ConvertToTime()
		err := store.SetBranchBrandGoalV2(ctx, generated.SetBranchBrandGoalV2Params{
			BranchID:      sql.NullInt32{Int32: branchId, Valid: true},
			BrandID:       sql.NullInt32{Int32: brandId, Valid: true},
			LeadMeasureID: leadMeasureID,
			Value:         goal,
			DateFrom:      from,
			DateTo:        to,
			Type:          string(entities.GoalTypeBranchBrand),
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	}
}
