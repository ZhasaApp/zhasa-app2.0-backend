package repository

import (
	"context"
	"time"
	entities2 "zhasa2.0/branch/entities"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/manager/entities"
	sale "zhasa2.0/sale/entities"
	"zhasa2.0/statistic"
	"zhasa2.0/statistic/repository"
)

/*
	SalesManagerRepository responsible to provide all data and data operations related to SalesManager
	Also contains Statistic interface, which gives all statistic related to SalesManager
*/
type SalesManagerRepository interface {
	SaveSale(salesDate time.Time, amount sale.SaleAmount, saleTypeId sale.SaleTypeId) error
	repository.StatisticRepository
	ProvideRankedSalesManagersList(from time.Time, to time.Time, size int32, page int32) (*entities.SalesManagers, error)
	ProvideBranchSalesManagers(id entities2.BranchId) (*entities.SalesManagers, error)
}

/*
	SalesManagerRepository implementation for real db data
*/
type PostgresSalesManagerRepository struct {
	sm      entities.SalesManager
	str     sale.SaleTypeRepository
	ctx     context.Context
	querier generated.Querier
}

func (p PostgresSalesManagerRepository) SaveSale(salesDate time.Time, amount sale.SaleAmount, saleTypeId sale.SaleTypeId) error {
	params := generated.AddSaleOrReplaceParams{
		SalesManagerID: int32(p.sm.Id),
		Date:           salesDate,
		Amount:         amount.Amount,
		SaleTypeID:     int32(saleTypeId),
	}

	return p.querier.AddSaleOrReplace(p.ctx, params)
}

func (p *PostgresSalesManagerRepository) ProvideSums(from time.Time, to time.Time) (*statistic.SaleSumByType, error) {
	arg := generated.GetSalesManagerSumsByTypeParams{
		Date:   from,
		Date_2: to,
		ID:     int32(p.sm.Id),
	}
	data, err := p.querier.GetSalesManagerSumsByType(p.ctx, arg)

	if err != nil {
		return nil, err
	}

	result := p.mapSalesSumsByType(data)
	return &result, err
}

func (p *PostgresSalesManagerRepository) mapSalesSumsByType(rows []generated.GetSalesManagerSumsByTypeRow) statistic.SaleSumByType {
	saleSumsByType := make(map[sale.SaleType]sale.SaleAmount)

	for _, row := range rows {
		saleAmount := sale.SaleAmount{Amount: row.TotalSales}
		saleType, err := p.str.GetSaleType(sale.SaleTypeId(row.SaleTypeID))
		if err != nil {
			return nil
		}
		saleSumsByType[*saleType] = saleAmount
	}

	return saleSumsByType
}

func (p *PostgresSalesManagerRepository) ProvideRankedSalesManagersList(from time.Time, to time.Time, size int32, page int32) (*entities.SalesManagers, error) {
	params := generated.GetRankedSalesManagersParams{
		Date:   from,
		Date_2: to,
		Limit:  size,
		Offset: page,
	}
	data, err := p.querier.GetRankedSalesManagers(p.ctx, params)
	if err != nil {
		return nil, err
	}
	var managers entities.SalesManagers
	for _, row := range data {
		salesManager := entities.SalesManager{
			Id:        entities.SalesManageId(row.SalesManagerID),
			FirstName: row.FirstName,
			LastName:  row.LastName,
			AvatarUrl: row.AvatarUrl,
		}
		managers = append(managers, salesManager)
	}
	return &managers, nil
}
