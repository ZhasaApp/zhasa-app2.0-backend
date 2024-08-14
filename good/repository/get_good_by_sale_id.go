package good

import (
	"context"
	"database/sql"
	generated "zhasa2.0/db/sqlc"
	good "zhasa2.0/good/entities"
)

type GetGoodBySaleIdFunc func(saleId int32) (*good.Good, error)

func NewGoodBySaleIdFunc(ctx context.Context, querier generated.Querier) GetGoodBySaleIdFunc {
	return func(saleId int32) (*good.Good, error) {
		res, err := querier.GetGoodBySaleId(ctx, saleId)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		return &good.Good{
			Id:          res.ID,
			Name:        res.Name,
			Description: res.Description,
		}, nil
	}
}
