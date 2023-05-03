package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/sale/entities"
)

type SaleTypeMap map[entities.SaleTypeId]*entities.SaleType

type SaleTypeRepository interface {
	GetSaleType(id entities.SaleTypeId) (*entities.SaleType, error)
	CreateSaleType(body entities.CreateSaleTypeBody) (entities.SaleTypeId, error)
}

type DBSaleTypeRepository struct {
	ctx     context.Context
	querier generated.Querier
	cache   SaleTypeMap
}

func NewSaleTypeRepository(ctx context.Context, querier generated.Querier) SaleTypeRepository {
	cache := make(SaleTypeMap)
	return DBSaleTypeRepository{
		ctx:     ctx,
		querier: querier,
		cache:   cache,
	}
}

func (str DBSaleTypeRepository) CreateSaleType(body entities.CreateSaleTypeBody) (entities.SaleTypeId, error) {
	params := generated.CreateSaleTypeParams{
		Title:       body.Title,
		Description: body.Description,
	}
	id, err := str.querier.CreateSaleType(str.ctx, params)

	if err != nil {
		return 0, err
	}
	return entities.SaleTypeId(id), nil
}

func (str DBSaleTypeRepository) GetSaleType(id entities.SaleTypeId) (*entities.SaleType, error) {
	saleType, found := str.cache[id]

	if found {
		return saleType, nil
	}

	saleTypeDb, err := str.querier.GetSaleTypeById(str.ctx, int32(id))
	if err != nil {
		return nil, err
	}

	newSaleType := &entities.SaleType{
		Id:          entities.SaleTypeId(saleTypeDb.ID),
		Title:       saleTypeDb.Title,
		Description: saleTypeDb.Description,
	}

	return newSaleType, nil
}
