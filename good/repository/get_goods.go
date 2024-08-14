package good

import (
	"context"
	"database/sql"
	generated "zhasa2.0/db/sqlc"
	good "zhasa2.0/good/entities"
)

type GetGoodsByBrandIdFunc func(brandId int32) ([]good.Good, error)

func NewGetGoodsByBrandIdFunc(ctx context.Context, querier generated.Querier) GetGoodsByBrandIdFunc {
	return func(brandId int32) ([]good.Good, error) {
		goods, err := querier.GetGoodsByBrandId(ctx, brandId)

		if err == sql.ErrNoRows {
			return []good.Good{}, nil
		}

		if err != nil {
			return nil, err
		}
		var result []good.Good
		for _, g := range goods {
			result = append(result, good.Good{
				Id:          g.ID,
				Name:        g.Name,
				Description: g.Description,
			})
		}
		return result, nil
	}
}
