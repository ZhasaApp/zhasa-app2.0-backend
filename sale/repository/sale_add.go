package repository

import (
	"context"
	"fmt"
	"time"
	generated "zhasa2.0/db/sqlc"
)

type SaleAddFunc func(saleAmount int64, saleDate time.Time, brandId int32, saleTypeId int32, saleDescription string) (int32, error)

func NewSaleAddFunc(ctx context.Context, store generated.SaleStore) SaleAddFunc {
	return func(saleAmount int64, saleDate time.Time, brandId int32, saleTypeId int32, saleDescription string) (int32, error) {
		sale, err := store.AddBrandSaleTx(ctx, generated.AddSaleOrReplaceParams{
			UserID:      1,
			SaleDate:    saleDate,
			Amount:      saleAmount,
			SaleTypeID:  saleTypeId,
			Description: saleDescription,
		}, brandId)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		return sale.ID, nil
	}
}
