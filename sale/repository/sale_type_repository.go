package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
	. "zhasa2.0/sale/entities"
	"zhasa2.0/statistic"
)

type SaleTypeMap map[SaleTypeId]*SaleType

type SaleTypeRepository interface {
	GetSaleType(id SaleTypeId) (*SaleType, error)
	CreateSaleType(body CreateSaleTypeBody) (SaleTypeId, error)
	MapSalesSumsByType(rows []SumsByTypeRow) statistic.SaleSumByType
}

type DBSaleTypeRepository struct {
	ctx     context.Context
	querier generated.Querier
	cache   SaleTypeMap
}

func (str DBSaleTypeRepository) MapSalesSumsByType(rows []SumsByTypeRow) statistic.SaleSumByType {
	saleSumsByType := make(map[SaleType]SaleAmount)

	for _, row := range rows {
		saleAmount := SaleAmount(row.TotalSales)
		saleType, err := str.GetSaleType(SaleTypeId(row.SaleTypeID))
		if err != nil {
			return nil
		}
		saleSumsByType[*saleType] = saleAmount
	}

	return saleSumsByType
}

func NewSaleTypeRepository(ctx context.Context, querier generated.Querier) SaleTypeRepository {
	cache := make(SaleTypeMap)
	return DBSaleTypeRepository{
		ctx:     ctx,
		querier: querier,
		cache:   cache,
	}
}

func (str DBSaleTypeRepository) CreateSaleType(body CreateSaleTypeBody) (SaleTypeId, error) {
	params := generated.CreateSaleTypeParams{
		Title:       body.Title,
		Description: body.Description,
	}
	id, err := str.querier.CreateSaleType(str.ctx, params)

	if err != nil {
		return 0, err
	}
	return SaleTypeId(id), nil
}

func (str DBSaleTypeRepository) GetSaleType(id SaleTypeId) (*SaleType, error) {
	saleType, found := str.cache[id]

	if found {
		return saleType, nil
	}

	saleTypeDb, err := str.querier.GetSaleTypeById(str.ctx, int32(id))
	if err != nil {
		return nil, err
	}

	newSaleType := &SaleType{
		Id:          SaleTypeId(saleTypeDb.ID),
		Title:       saleTypeDb.Title,
		Description: saleTypeDb.Description,
	}

	return newSaleType, nil
}
