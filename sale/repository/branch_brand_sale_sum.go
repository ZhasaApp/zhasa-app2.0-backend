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
		row, err := store.GetSaleSumByBranchByTypeByBrand(ctxt, generated.GetSaleSumByBranchByTypeByBrandParams{
			SaleDate:   from,
			SaleDate_2: to,
			ID:         branchId,
			ID_2:       brandId,
			ID_3:       saleTypeId,
		})
		if err == sql.ErrNoRows {
			return 0, nil
		}
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		return GetInt64FromInterface(row.TotalSales), nil
	}
}

func GetInt64FromInterface(val interface{}) int64 {
	switch v := val.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	default:
		return 0
	}
}
