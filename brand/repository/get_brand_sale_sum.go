package repository

import (
	"context"
	"database/sql"
	"fmt"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/statistic"
)

type GetBrandSaleSumFunc func(brandId int32, saleTypeId int32, period statistic.Period) (int64, error)

func NewGetBrandSaleSumFunc(ctxt context.Context, store generated.SaleStore) GetBrandSaleSumFunc {
	return func(brandId int32, saleTypeId int32, period statistic.Period) (int64, error) {
		from, to := period.ConvertToTime()
		row, err := store.GetBrandSaleSumByGivenDateRange(ctxt, generated.GetBrandSaleSumByGivenDateRangeParams{
			BrandID:    brandId,
			SaleTypeID: saleTypeId,
			SaleDate:   from,
			SaleDate_2: to,
		})
		if err == sql.ErrNoRows {
			return 0, nil
		}
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		return row, nil
	}
}
