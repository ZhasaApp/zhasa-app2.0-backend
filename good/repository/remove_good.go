package good

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type DeleteGoodFunc func(goodId int32) error

func NewDeleteGoodFunc(ctx context.Context, querier generated.Querier) DeleteGoodFunc {
	return func(goodId int32) error {
		err := querier.DisableGood(ctx, goodId)
		return err
	}
}
