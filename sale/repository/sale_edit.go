package repository

import (
	"context"
	"fmt"
	"time"
	generated "zhasa2.0/db/sqlc"
)

type SaleEditFunc func(saleAmount int64, saleDate time.Time, saleId, userId, saleTypeId int32, saleDescription string) (int32, error)

func NewSaleEditFunc(ctx context.Context, store generated.SaleStore) func(saleAmount int64, saleDate time.Time, saleId int32, userId int32, saleTypeId int32, saleDescription string) (int32, error) {
	return func(saleAmount int64, saleDate time.Time, saleId, userId, saleTypeId int32, saleDescription string) (int32, error) {
		sale, err := store.EditSale(ctx, generated.EditSaleParams{
			ID:          saleId,
			UserID:      userId,
			SaleDate:    saleDate,
			Amount:      saleAmount,
			SaleTypeID:  saleTypeId,
			Description: saleDescription,
		})
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		return sale.ID, nil
	}
}
