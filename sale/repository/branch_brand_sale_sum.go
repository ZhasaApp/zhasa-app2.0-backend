package repository

import (
	"context"
	"database/sql"
	"fmt"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/statistic"
)

type GetBranchBrandSaleSumFunc func(branchId int32, brandId int32, saleTypeId int32, period statistic.Period) (int64, error)

func NewGetBranchBrandSaleSumFunc(ctxt context.Context, store generated.SaleStore) GetBranchBrandSaleSumFunc {
	return func(branchId int32, brandId int32, saleTypeId int32, period statistic.Period) (int64, error) {
		from, to := period.ConvertToTime()
		row, err := store.GetBranchBrandSaleSumByGivenDateRange(ctxt, generated.GetBranchBrandSaleSumByGivenDateRangeParams{
			BranchID:   branchId,
			BrandID:    brandId,
			SaleTypeID: saleTypeId,
			SaleDate:   from,
			SaleDate_2: to,
		})
		fmt.Println("row: ", getInt64FromInterface(row))
		if err == sql.ErrNoRows {
			return 0, nil
		}
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		return getInt64FromInterface(row), nil
	}
}

func getInt64FromInterface(val interface{}) int64 {
	switch v := val.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case int32:
		return int64(v) // Convert int32 to int64
	case float64:
		return int64(v)
	default:
		return 0
	}
}
