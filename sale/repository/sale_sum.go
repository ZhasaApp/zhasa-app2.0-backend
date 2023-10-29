package repository

import (
	"context"
	"database/sql"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/statistic"
)

type GetSaleSumByUserBrandTypePeriodFunc func(userId int32, brandId int32, saleTypeId int32, period statistic.Period) (int64, error)

func NewGetSaleSumByUserBrandTypePeriodFunc(ctx context.Context, store generated.SaleStore) GetSaleSumByUserBrandTypePeriodFunc {
	return func(userId int32, brandId int32, saleTypeId int32, period statistic.Period) (int64, error) {
		from, to := period.ConvertToTime()
		sum, err := store.GetSumByUserIdBrandIdPeriodSaleTypeId(ctx, generated.GetSumByUserIdBrandIdPeriodSaleTypeIdParams{
			UserID:     userId,
			BrandID:    brandId,
			SaleTypeID: saleTypeId,
			SaleDate:   from,
			SaleDate_2: to,
		})
		if err == sql.ErrNoRows {
			return 0, nil
		}
		if err != nil {
			return 0, err
		}
		return sum, nil
	}
}
