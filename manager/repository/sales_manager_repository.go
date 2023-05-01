package repository

import (
	"context"
	"time"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/manager/entities"
	sale "zhasa2.0/sale/entities"
	repository2 "zhasa2.0/sale/repository"
	"zhasa2.0/statistic"
	"zhasa2.0/statistic/repository"
)

/*
	SalesManagerRepository responsible to provide all data and data operations related to SalesManager
	Also contains Statistic interface, which gives all statistic related to SalesManager
*/
type SalesManagerRepository interface {
	CreateSalesManager(userId int32, branchId int32) error
	SaveSale(salesManagerId entities.SalesManagerId, salesDate time.Time, amount sale.SaleAmount, saleTypeId sale.SaleTypeId) error
	repository.StatisticRepository
	ProvideRankedSalesManagersList(from time.Time, to time.Time, size int32, page int32) (*entities.SalesManagers, error)
	GetSalesManagerByUserId(userId int32) (*entities.SalesManager, error)
}

/*
	SalesManagerRepository implementation for real db data
*/
type PostgresSalesManagerRepository struct {
	repository2.SaleTypeRepository
	ctx     context.Context
	querier generated.Querier
}

func NewSalesManagerRepository(typeRepository repository2.SaleTypeRepository, ctx context.Context, querier generated.Querier) SalesManagerRepository {
	return PostgresSalesManagerRepository{
		SaleTypeRepository: typeRepository,
		ctx:                ctx,
		querier:            querier,
	}
}

func (p PostgresSalesManagerRepository) CreateSalesManager(userId int32, branchId int32) error {
	params := generated.CreateSalesManagerParams{
		UserID:   userId,
		BranchID: branchId,
	}
	return p.querier.CreateSalesManager(p.ctx, params)
}

func (p PostgresSalesManagerRepository) SaveSale(salesManagerId entities.SalesManagerId, salesDate time.Time, amount sale.SaleAmount, saleTypeId sale.SaleTypeId) error {
	params := generated.AddSaleOrReplaceParams{
		SalesManagerID: int32(salesManagerId),
		SaleDate:       salesDate,
		Amount:         int64(amount),
		SaleTypeID:     int32(saleTypeId),
	}

	return p.querier.AddSaleOrReplace(p.ctx, params)
}

func (p PostgresSalesManagerRepository) ProvideSums(salesManagerId entities.SalesManagerId, from time.Time, to time.Time) (*statistic.SaleSumByType, error) {
	arg := generated.GetSalesManagerSumsByTypeParams{
		SaleDate:   from,
		SaleDate_2: to,
		ID:         int32(salesManagerId),
	}
	data, err := p.querier.GetSalesManagerSumsByType(p.ctx, arg)

	if err != nil {
		return nil, err
	}

	result := p.mapSalesSumsByType(data)
	return &result, err
}

func (p PostgresSalesManagerRepository) mapSalesSumsByType(rows []generated.GetSalesManagerSumsByTypeRow) statistic.SaleSumByType {
	saleSumsByType := make(map[sale.SaleType]sale.SaleAmount)

	for _, row := range rows {
		saleAmount := sale.SaleAmount(row.TotalSales)
		saleType, err := p.GetSaleType(sale.SaleTypeId(row.SaleTypeID))
		if err != nil {
			return nil
		}
		saleSumsByType[*saleType] = saleAmount
	}

	return saleSumsByType
}

func (p PostgresSalesManagerRepository) ProvideRankedSalesManagersList(from time.Time, to time.Time, size int32, page int32) (*entities.SalesManagers, error) {
	params := generated.GetRankedSalesManagersParams{
		SaleDate:   from,
		SaleDate_2: to,
		Limit:      size,
		Offset:     page,
	}
	data, err := p.querier.GetRankedSalesManagers(p.ctx, params)
	if err != nil {
		return nil, err
	}
	var managers entities.SalesManagers
	for _, row := range data {
		salesManager := entities.SalesManager{
			Id:        entities.SalesManagerId(row.SalesManagerID),
			FirstName: row.FirstName,
			LastName:  row.LastName,
			AvatarUrl: row.AvatarUrl,
		}
		managers = append(managers, salesManager)
	}
	return &managers, nil
}

func (p PostgresSalesManagerRepository) GetSalesManagerByUserId(userId int32) (*entities.SalesManager, error) {
	data, err := p.querier.GetSalesManagerByUserId(p.ctx, userId)
	if err != nil {
		return nil, err
	}

	salesManager := entities.SalesManager{
		Id:        entities.SalesManagerId(data.SalesManagerID),
		FirstName: data.FirstName,
		LastName:  data.LastName,
		AvatarUrl: data.AvatarUrl,
	}
	return &salesManager, err
}
