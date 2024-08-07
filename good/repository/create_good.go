package good

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type CreateGoodFunc func(name, description string) (int32, error)

func NewCreateGoodFunc(ctx context.Context, querier generated.Querier) CreateGoodFunc {
	return func(name, description string) (int32, error) {
		id, err := querier.CreateGood(ctx, generated.CreateGoodParams{
			Name:        name,
			Description: description,
		})
		if err != nil {
			return 0, err
		}
		return id, nil
	}
}
