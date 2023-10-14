package repository

import (
	"context"
	"fmt"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/statistic"
)

type SetBranchBrandSaleTypeGoal func(branchBrand int32, saleTypeId int32, goal int64, period statistic.Period) error

func NewSetBranchGoalFunc(ctx context.Context, store generated.BranchStore) SetBranchBrandSaleTypeGoal {
	return func(branchBrand int32, saleTypeId int32, goal int64, period statistic.Period) error {
		from, to := period.ConvertToTime()
		err := store.SetBranchBrandGoal(ctx, generated.SetBranchBrandGoalParams{
			BranchBrand: branchBrand,
			SaleTypeID:  saleTypeId,
			Value:       goal,
			FromDate:    from,
			ToDate:      to,
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	}
}
