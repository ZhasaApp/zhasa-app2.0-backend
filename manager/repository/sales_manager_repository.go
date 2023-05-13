package repository

import (
	"context"
	"database/sql"
	"time"
	generated "zhasa2.0/db/sqlc"
	. "zhasa2.0/manager/entities"
	. "zhasa2.0/sale/entities"
	repository2 "zhasa2.0/sale/repository"
	"zhasa2.0/statistic"
	. "zhasa2.0/statistic/entities"
	"zhasa2.0/statistic/repository"
)

/*
	SalesManagerRepository responsible to provide all data and data operations related to SalesManager
	Also contains Statistic interface, which gives all statistic related to SalesManager
*/
type SalesManagerRepository interface {
	CreateSalesManager(userId int32, branchId int32) error
	SaveSale(salesManagerId SalesManagerId, salesDate time.Time, amount SaleAmount, saleTypeId SaleTypeId) error
	repository.StatisticRepository
	ProvideRankedSalesManagersList(from time.Time, to time.Time, size int32, page int32) (*SalesManagers, error)
	GetSalesManagerByUserId(userId int32) (*SalesManager, error)
	GetSalesManagerGoalAmount(salesManagerId SalesManagerId, from time.Time, to time.Time) (SaleAmount, error)
	GetMonthlyYearSaleStatistic(saleManagerId SalesManagerId, year int32) (*[]YearStatisticByMonth, error)
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

func (p PostgresSalesManagerRepository) SaveSale(salesManagerId SalesManagerId, salesDate time.Time, amount SaleAmount, saleTypeId SaleTypeId) error {
	params := generated.AddSaleOrReplaceParams{
		SalesManagerID: int32(salesManagerId),
		SaleDate:       salesDate,
		Amount:         int64(amount),
		SaleTypeID:     int32(saleTypeId),
	}

	return p.querier.AddSaleOrReplace(p.ctx, params)
}

func (p PostgresSalesManagerRepository) ProvideSums(salesManagerId SalesManagerId, from time.Time, to time.Time) (*statistic.SaleSumByType, error) {
	arg := generated.GetSalesManagerSumsByTypeParams{
		SaleDate:       from,
		SaleDate_2:     to,
		SalesManagerID: int32(salesManagerId),
	}
	data, err := p.querier.GetSalesManagerSumsByType(p.ctx, arg)

	if err != nil {
		return nil, err
	}

	result := p.mapSalesSumsByType(data)
	return &result, err
}

func (p PostgresSalesManagerRepository) mapSalesSumsByType(rows []generated.GetSalesManagerSumsByTypeRow) statistic.SaleSumByType {
	saleSumsByType := make(map[SaleType]SaleAmount)

	for _, row := range rows {
		saleAmount := SaleAmount(row.TotalSales)
		saleType, err := p.GetSaleType(SaleTypeId(row.SaleTypeID))
		if err != nil {
			return nil
		}
		saleSumsByType[*saleType] = saleAmount
	}

	return saleSumsByType
}

func (p PostgresSalesManagerRepository) ProvideRankedSalesManagersList(from time.Time, to time.Time, size int32, page int32) (*SalesManagers, error) {
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
	var managers SalesManagers
	for _, row := range data {
		salesManager := SalesManager{
			Id:        SalesManagerId(row.SalesManagerID),
			FirstName: row.FirstName,
			LastName:  row.LastName,
			AvatarUrl: row.AvatarUrl,
		}
		managers = append(managers, salesManager)
	}
	return &managers, nil
}

func (p PostgresSalesManagerRepository) GetSalesManagerByUserId(userId int32) (*SalesManager, error) {
	data, err := p.querier.GetSalesManagerByUserId(p.ctx, userId)
	if err != nil {
		return nil, err
	}

	salesManager := SalesManager{
		Id:        SalesManagerId(data.SalesManagerID),
		FirstName: data.FirstName,
		LastName:  data.LastName,
		AvatarUrl: data.AvatarUrl,
	}
	return &salesManager, err
}

func (p PostgresSalesManagerRepository) GetSalesManagerGoalAmount(smId SalesManagerId, from time.Time, to time.Time) (SaleAmount, error) {
	arg := generated.GetSalesManagerGoalByGivenDateRangeParams{
		SalesManagerID: int32(smId),
		FromDate:       from,
		ToDate:         to,
	}
	data, err := p.querier.GetSalesManagerGoalByGivenDateRange(p.ctx, arg)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return SaleAmount(data), nil
}

func (p PostgresSalesManagerRepository) GetMonthlyYearSaleStatistic(saleManagerId SalesManagerId, year int32) (*[]YearStatisticByMonth, error) {
	params := generated.GetSalesManagerYearStatisticParams{
		SalesManagerID: int32(saleManagerId),
		SaleDate:       year,
	}
	data, err := p.querier.GetSalesManagerYearStatistic(p.ctx, params)

	result := make([]YearStatisticByMonth, 0)

	if err == sql.ErrNoRows {
		return &result, nil
	}
	if err != nil {
		return nil, err
	}
	for _, row := range data {
		saleType, err := p.GetSaleType(SaleTypeId(row.SaleType))
		if err != nil {
			return nil, err
		}
		stat := YearStatisticByMonth{
			SaleType: *saleType,
			Month:    MonthNumber(row.MonthNumber),
			Amount:   SaleAmount(row.TotalAmount),
		}
		result = append(result, stat)
	}
	return &result, nil
}
