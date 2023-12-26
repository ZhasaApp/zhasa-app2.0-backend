package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/statistic"
)

type SetBrandSaleTypeGoalFunc func(brandId, saleTypeId int32, goal int64, period statistic.Period) error

func NewSetBrandSaleTypeGoalFunc(
	ctx context.Context,
	store generated.BranchStore,
) SetBrandSaleTypeGoalFunc {
	return func(brandId, saleTypeId int32, goal int64, period statistic.Period) error {
		from, to := period.ConvertToTime()
		err := store.SetBrandSaleTypeGoal(ctx, generated.SetBrandSaleTypeGoalParams{
			BrandID:    brandId,
			SaleTypeID: saleTypeId,
			Value:      goal,
			FromDate:   from,
			ToDate:     to,
		})
		if err != nil {
			return err
		}
		return nil
	}
}
