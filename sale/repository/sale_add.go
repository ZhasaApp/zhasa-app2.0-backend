package repository

import (
	"context"
	"fmt"
	"time"
	generated "zhasa2.0/db/sqlc"
)

type SaleAddBody struct {
	SaleAmount      int64     `json:"sale_amount"`
	SaleDate        time.Time `json:"sale_date"`
	UserId          int32     `json:"user_id"`
	BrandId         int32     `json:"brand_id"`
	SaleTypeId      int32     `json:"sale_type_id"`
	SaleDescription string    `json:"sale_description"`
}

type SaleAddFunc func(body SaleAddBody) (int32, error)

func NewSaleAddFunc(ctx context.Context, store generated.SaleStore) SaleAddFunc {
	return func(body SaleAddBody) (int32, error) {
		sale, err := store.AddBrandSaleTx(ctx, generated.AddSaleOrReplaceParams{
			UserID:      body.UserId,
			SaleDate:    body.SaleDate,
			Amount:      body.SaleAmount,
			SaleTypeID:  body.SaleTypeId,
			Description: body.SaleDescription,
		}, body.BrandId)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		return sale.ID, nil
	}
}
