package entities

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type SaleTypeId int32

type SaleType struct {
	Id          SaleTypeId
	Title       string
	Description string
}

type SaleTypeMap map[SaleTypeId]*SaleType

type SaleTypeRepository interface {
	GetSaleType(id SaleTypeId) (*SaleType, error)
}

type DBSaleTypeRepository struct {
	ctx     context.Context
	querier generated.Querier
	cache   SaleTypeMap
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
		Description: saleTypeDb.Description.String,
	}

	return newSaleType, nil
}
