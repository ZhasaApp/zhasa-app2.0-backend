package good

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type AddGoodToBrandFunc func(goodId, brandId int32) error

func NewAddGoodToBrandFunc(ctx context.Context, querier generated.Querier) AddGoodToBrandFunc {
	return func(goodId, brandId int32) error {
		err := querier.AddGoodToBrand(ctx, generated.AddGoodToBrandParams{
			BrandID: brandId,
			GoodID:  goodId,
		})
		if err != nil {
			return err
		}
		return nil
	}
}
