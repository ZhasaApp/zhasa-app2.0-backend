package repository

import (
	"context"
	"database/sql"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/statistic"
)

type SetUserBrandGoalFunc func(userId, brandId, leadMeasureID int32, goal int64, period statistic.Period) error

func NewSetUserBrandGoalFunc(ctx context.Context, store generated.UserBrandStore) SetUserBrandSaleTypeGoalFunc {
	return func(userId, brandId, leadMeasureID int32, goal int64, period statistic.Period) error {
		from, to := period.ConvertToTime()
		err := store.SetUserBrandGoalV2(ctx, generated.SetUserBrandGoalV2Params{
			UserID:        sql.NullInt32{Int32: userId, Valid: true},
			BrandID:       sql.NullInt32{Int32: brandId, Valid: true},
			LeadMeasureID: leadMeasureID,
			Value:         goal,
			DateFrom:      from,
			DateTo:        to,
		})
		if err != nil {
			return err
		}
		return nil
	}
}
