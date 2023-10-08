package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
	. "zhasa2.0/sale/entities"
	"zhasa2.0/statistic"
)

type SaleTypeMap map[int32]*SaleType

type SaleTypeRepository interface {
	GetSaleType(id int32) (*SaleType, error)
	CreateSaleType(body CreateSaleTypeBody) (int32, error)
	MapSalesSumsByType(rows []SumsByTypeRow) statistic.SaleSumByType
	GetSaleTypes() (*[]SaleType, error)
}

type DBSaleTypeRepository struct {
	ctx     context.Context
	querier generated.Querier
	cache   SaleTypeMap
}

func (str DBSaleTypeRepository) GetSaleTypes() (*[]SaleType, error) {
	rows, err := str.querier.GetSalesTypes(str.ctx)

	if err != nil {
		return nil, err
	}

	saleTypes := make([]SaleType, 0)

	for _, row := range rows {
		saleTypes = append(saleTypes, SaleType{
			Id:          row.ID,
			Title:       row.Title,
			Description: row.Description,
			Color:       row.Color,
			Gravity:     row.Gravity,
			ValueType:   string(row.ValueType),
		})
	}

	return &saleTypes, nil
}

func (str DBSaleTypeRepository) MapSalesSumsByType(rows []SumsByTypeRow) statistic.SaleSumByType {
	saleSumsByType := make(map[SaleType]SaleAmount)

	for _, row := range rows {
		saleAmount := SaleAmount(row.TotalSales)
		saleType, err := str.GetSaleType(row.SaleTypeID)
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

func (str DBSaleTypeRepository) CreateSaleType(body CreateSaleTypeBody) (int32, error) {
	params := generated.CreateSaleTypeParams{
		Title:       body.Title,
		Description: body.Description,
	}
	id, err := str.querier.CreateSaleType(str.ctx, params)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (str DBSaleTypeRepository) GetSaleType(id int32) (*SaleType, error) {
	saleType, found := str.cache[id]

	if found {
		return saleType, nil
	}

	saleTypeDb, err := str.querier.GetSaleTypeById(str.ctx, id)
	if err != nil {
		return nil, err
	}

	newSaleType := &SaleType{
		Id:          saleTypeDb.ID,
		Title:       saleTypeDb.Title,
		Description: saleTypeDb.Description,
		Color:       saleTypeDb.Color,
		ValueType:   string(saleTypeDb.ValueType),
	}

	return newSaleType, nil
}
