package repository

import (
	"context"
	"fmt"
	. "zhasa2.0/db/sqlc"
)

type SaleRepository interface {
	AddOrEdit(saleToCreate AddSaleOrReplaceParams, brandId int32) error
	GetSumByUserIdBrandIdPeriodSaleTypeId(params GetSaleSumByUserIdBrandIdPeriodSaleTypeIdParams) (int64, error)
}

type DBSaleRepository struct {
	ctx   context.Context
	store SaleStore
}

func (d DBSaleRepository) AddOrEdit(saleToCreate AddSaleOrReplaceParams, brandId int32) error {
	_, err := d.store.AddBrandSaleTx(d.ctx, saleToCreate, brandId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (d DBSaleRepository) GetSumByUserIdBrandIdPeriodSaleTypeId(params GetSaleSumByUserIdBrandIdPeriodSaleTypeIdParams) (int64, error) {
	amount, err := d.store.GetSaleSumByUserIdBrandIdPeriodSaleTypeId(d.ctx, params)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return amount, err
}

func NewSaleRepo(ctx context.Context, store *DBStore) SaleRepository {
	return DBSaleRepository{
		ctx:   ctx,
		store: store,
	}
}
