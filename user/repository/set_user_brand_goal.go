package repository

import (
	"context"
	"fmt"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/statistic"
)

type SetUserBrandSaleTypeGoalFunc func(userBrand int32, saleTypeId int32, goal int64, period statistic.Period) error

func NewSetUserBrandSaleTypeGoalFunc(ctx context.Context, store generated.UserBrandStore) SetUserBrandSaleTypeGoalFunc {
	return func(userBrand int32, saleTypeId int32, goal int64, period statistic.Period) error {
		from, to := period.ConvertToTime()
		err := store.SetUserBrandGoal(ctx, generated.SetUserBrandGoalParams{
			UserBrand:  userBrand,
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
