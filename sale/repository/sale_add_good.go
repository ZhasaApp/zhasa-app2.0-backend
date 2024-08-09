package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type AddGoodToSaleBody struct {
	SaleId int32 `json:"sale_id"`
	GoodId int32 `json:"good_id"`
}

type SaleAddWithGoodFunc func(body AddGoodToSaleBody) error

func NewSaleAddWithGoodFunc(ctx context.Context, store generated.SaleStore) SaleAddWithGoodFunc {
	return func(body AddGoodToSaleBody) error {
		return store.AddGoodToSaleTx(ctx, body.SaleId, body.GoodId)
	}
}
