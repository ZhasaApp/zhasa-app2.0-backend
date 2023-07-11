package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
	. "zhasa2.0/base"
	. "zhasa2.0/branch/entities"
	. "zhasa2.0/db/hand-made"
	generated "zhasa2.0/db/sqlc"
	. "zhasa2.0/manager/entities"
	. "zhasa2.0/sale/entities"
	repository2 "zhasa2.0/sale/repository"
	. "zhasa2.0/statistic/entities"
	. "zhasa2.0/user/entities"
)

/*
	SalesManagerRepository responsible to provide all data and data operations related to SalesManager
	Also contains Statistic interface, which gives all statistic related to SalesManager
*/
type SalesManagerRepository interface {
	CreateSalesManager(userId int32, branchId int32) error
	SaveSale(salesManagerId SalesManagerId, salesDate time.Time, amount SaleAmount, saleTypeId SaleTypeId, description SaleDescription) (*Sale, error)
	GetSalesManagerByUserId(userId int32) (*SalesManager, error)
	GetMonthlyYearSaleStatistic(smId SalesManagerId, year int32) (*[]MonthlyYearStatistic, error)
	GetManagerSales(salesManagerId SalesManagerId, pagination Pagination) (*[]Sale, error)
	GetManagerSalesByPeriod(salesManagerId SalesManagerId, pagination Pagination, from time.Time, to time.Time) (*[]Sale, error)
	GetSalesManagerSalesCount(salesManagerId SalesManagerId) (int32, error)
	GetSalesManagersListOrderedByRatio(pagination Pagination, from time.Time, to time.Time) (*[]SalesManager, error)
	DeleteSale(saleId SaleId) (*Sale, error)
	EditSale(body EditSaleBody) (*Sale, error)
}

type SalesManagerStatisticRepository interface {
	GetSalesSumBySaleTypeAndManager(smId SalesManagerId, typeId SaleTypeId, from, to time.Time) (SaleAmount, error)
	GetSalesGoalBySaleTypeAndManager(smId SalesManagerId, typeId SaleTypeId, from, to time.Time) (SaleAmount, error)
	GetSalesManagerRatioByPeriod(smId SalesManagerId, from, to time.Time) (Percent, error)
	SetRatioByPeriod(smId SalesManagerId, ratio Percent, from, to time.Time) error
}

type PostgresSalesManagerStatisticRepository struct {
	repository2.SaleTypeRepository
	ctx     context.Context
	querier generated.Querier
}

func (p PostgresSalesManagerStatisticRepository) SetRatioByPeriod(smId SalesManagerId, ratio Percent, from, to time.Time) error {
	err := p.querier.SetSMRatio(p.ctx, generated.SetSMRatioParams{
		FromDate:       from,
		ToDate:         to,
		Ratio:          float64(ratio),
		SalesManagerID: int32(smId),
	})
	return err
}

func (p PostgresSalesManagerStatisticRepository) GetSalesManagerRatioByPeriod(smId SalesManagerId, from, to time.Time) (Percent, error) {
	ratio, err := p.querier.GetSMRatio(p.ctx, generated.GetSMRatioParams{
		FromDate:       from,
		ToDate:         to,
		SalesManagerID: int32(smId),
	})

	if err == sql.ErrNoRows {
		return Percent(0), nil
	}
	if err != nil {
		return Percent(0), err
	}

	return Percent(ratio), nil
}

/*
	SalesManagerRepository implementation for real db data
*/
type PostgresSalesManagerRepository struct {
	repository2.SaleTypeRepository
	ctx           context.Context
	querier       generated.Querier
	customQuerier CustomQuerier
}

func (p PostgresSalesManagerRepository) EditSale(body EditSaleBody) (*Sale, error) {
	params := generated.EditSaleByIdParams{
		ID:          body.ID,
		SaleTypeID:  body.TypeID,
		SaleDate:    body.Date,
		Amount:      body.Value,
		Description: body.Title,
	}

	row, err := p.querier.EditSaleById(p.ctx, params)

	if err != nil {
		return nil, err
	}

	return &Sale{
		Id:              SaleId(row.ID),
		SaleManagerId:   SalesManagerId(row.SalesManagerID),
		SaleType:        SaleType{},
		SalesAmount:     SaleAmount(row.Amount),
		SaleDate:        row.SaleDate,
		SaleDescription: SaleDescription(row.Description),
	}, nil
}

func (p PostgresSalesManagerRepository) DeleteSale(saleId SaleId) (*Sale, error) {
	row, err := p.querier.DeleteSaleById(p.ctx, int32(saleId))
	if err != nil {
		return nil, err
	}
	return &Sale{
		Id:              SaleId(row.ID),
		SaleManagerId:   SalesManagerId(row.SalesManagerID),
		SaleType:        SaleType{},
		SalesAmount:     0,
		SaleDate:        row.SaleDate,
		SaleDescription: "",
	}, nil
}

func (p PostgresSalesManagerRepository) GetSalesManagersListOrderedByRatio(pagination Pagination, from time.Time, to time.Time) (*[]SalesManager, error) {
	data, err := p.querier.GetOrderedSalesManagers(p.ctx, generated.GetOrderedSalesManagersParams{
		FromDate: from,
		ToDate:   to,
		Limit:    pagination.PageSize,
		Offset:   pagination.Page,
	})

	result := make([]SalesManager, 0)

	if err == sql.ErrNoRows {
		return &result, nil
	}
	if err != nil {
		return nil, err
	}

	for index, item := range data {
		result = append(result, SalesManager{
			Id:        SalesManagerId(item.SalesManagerID),
			FirstName: item.FirstName,
			LastName:  item.LastName,
			AvatarUrl: item.AvatarUrl,
			Branch: Branch{
				BranchId:    BranchId(item.BranchID),
				Title:       BranchTitle(item.BranchTitle),
				Description: "",
				Key:         "",
			},
			Ratio:       Percent(item.Ratio).GetRounded(),
			RatingPlace: RatingPlace((pagination.Page)*pagination.PageSize + int32(index) + int32(1)),
			UserId:      UserId(item.UserID),
		})
	}
	return &result, err
}

func (p PostgresSalesManagerRepository) GetManagerSalesByPeriod(salesManagerId SalesManagerId, pagination Pagination, from time.Time, to time.Time) (*[]Sale, error) {
	params := generated.GetManagerSalesByPeriodParams{
		SalesManagerID: int32(salesManagerId),
		SaleDate:       from,
		SaleDate_2:     to,
		Limit:          pagination.PageSize,
		Offset:         pagination.Page,
	}

	rows, err := p.querier.GetManagerSalesByPeriod(p.ctx, params)

	result := make([]Sale, 0)
	if err == sql.ErrNoRows {
		return &result, nil
	}
	if err != nil {
		return nil, err
	}

	for _, item := range rows {
		t, err := p.SaleTypeRepository.GetSaleType(SaleTypeId(item.SaleTypeID))
		if err != nil {
			return nil, err
		}
		result = append(result, Sale{
			Id:              SaleId(item.ID),
			SaleManagerId:   salesManagerId,
			SaleType:        *t,
			SalesAmount:     SaleAmount(item.Amount),
			SaleDate:        item.SaleDate,
			SaleDescription: SaleDescription(item.Description),
		})
	}
	return &result, err
}

func (p PostgresSalesManagerRepository) GetSalesManagerSalesCount(salesManagerId SalesManagerId) (int32, error) {
	count, err := p.querier.GetSalesCount(p.ctx, int32(salesManagerId))
	return int32(count), err
}

func (p PostgresSalesManagerRepository) GetManagerSales(salesManagerId SalesManagerId, pagination Pagination) (*[]Sale, error) {
	params := generated.GetManagerSalesParams{
		SalesManagerID: int32(salesManagerId),
		Limit:          pagination.PageSize,
		Offset:         pagination.Page,
	}
	data, err := p.querier.GetManagerSales(p.ctx, params)

	result := make([]Sale, 0)
	if err == sql.ErrNoRows {
		return &result, nil
	}
	if err != nil {
		return nil, err
	}

	for _, item := range data {
		t, err := p.SaleTypeRepository.GetSaleType(SaleTypeId(item.SaleTypeID))
		if err != nil {
			return nil, err
		}
		result = append(result, Sale{
			Id:              SaleId(item.ID),
			SaleManagerId:   salesManagerId,
			SaleType:        *t,
			SalesAmount:     SaleAmount(item.Amount),
			SaleDate:        item.SaleDate,
			SaleDescription: SaleDescription(item.Description),
		})
	}
	return &result, err
}

func NewSalesManagerRepository(typeRepository repository2.SaleTypeRepository, ctx context.Context, querier generated.Querier, customQuerier CustomQuerier) SalesManagerRepository {
	return PostgresSalesManagerRepository{
		SaleTypeRepository: typeRepository,
		ctx:                ctx,
		querier:            querier,
		customQuerier:      customQuerier,
	}
}

func NewSalesManagerStatisticRepository(typeRepository repository2.SaleTypeRepository, ctx context.Context, querier generated.Querier) SalesManagerStatisticRepository {
	return PostgresSalesManagerStatisticRepository{
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

func (p PostgresSalesManagerRepository) SaveSale(salesManagerId SalesManagerId, salesDate time.Time, amount SaleAmount, saleTypeId SaleTypeId, description SaleDescription) (*Sale, error) {
	params := generated.AddSaleOrReplaceParams{
		SalesManagerID: int32(salesManagerId),
		SaleDate:       salesDate,
		Amount:         int64(amount),
		SaleTypeID:     int32(saleTypeId),
		Description:    string(description),
	}

	row, err := p.querier.AddSaleOrReplace(p.ctx, params)
	if err != nil {
		return nil, err
	}

	saleType, err := p.GetSaleType(SaleTypeId(row.SaleTypeID))
	return &Sale{
		Id:              SaleId(row.ID),
		SaleManagerId:   SalesManagerId(row.SalesManagerID),
		SaleType:        *saleType,
		SalesAmount:     SaleAmount(row.Amount),
		SaleDate:        row.SaleDate,
		SaleDescription: SaleDescription(row.Description),
	}, nil
}

func (p PostgresSalesManagerRepository) GetSalesManagerByUserId(userId int32) (*SalesManager, error) {
	data, err := p.querier.GetSalesManagerByUserId(p.ctx, userId)
	if err == sql.ErrNoRows {
		return nil, errors.New("sales manager not found")
	}
	if err != nil {
		return nil, err
	}
	salesManager := SalesManager{
		Id:        SalesManagerId(data.SalesManagerID),
		FirstName: data.FirstName,
		LastName:  data.LastName,
		AvatarUrl: data.AvatarUrl,
		Branch: Branch{
			BranchId:    BranchId(data.BranchID),
			Title:       BranchTitle(data.BranchTitle),
			Description: "",
			Key:         "",
		},
	}
	return &salesManager, err
}

func (p PostgresSalesManagerStatisticRepository) GetSalesSumBySaleTypeAndManager(smId SalesManagerId, typeId SaleTypeId, from, to time.Time) (SaleAmount, error) {
	arg := generated.GetSalesManagerSumsByTypeParams{
		SalesManagerID: int32(smId),
		SaleDate:       from,
		SaleDate_2:     to,
		ID:             int32(typeId),
	}
	data, err := p.querier.GetSalesManagerSumsByType(p.ctx, arg)

	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return SaleAmount(data.TotalSales), nil
}

func (p PostgresSalesManagerStatisticRepository) GetSalesGoalBySaleTypeAndManager(smId SalesManagerId, typeId SaleTypeId, from, to time.Time) (SaleAmount, error) {
	arg := generated.GetSalesManagerGoalByGivenDateRangeAndSaleTypeParams{
		SalesManagerID: int32(smId),
		FromDate:       from,
		ToDate:         to,
		TypeID:         int32(typeId),
	}
	data, err := p.querier.GetSalesManagerGoalByGivenDateRangeAndSaleType(p.ctx, arg)

	if err != nil {
		return 0, err
	}

	return SaleAmount(data), nil
}

func (p PostgresSalesManagerRepository) GetMonthlyYearSaleStatistic(smId SalesManagerId, year int32) (*[]MonthlyYearStatistic, error) {

	result := make([]MonthlyYearStatistic, 0)

	saleTypes, err := p.GetSaleTypes()
	if err != nil {
		return nil, errors.New("error while getting sale types")
	}

	for _, saleType := range *saleTypes {
		for month := 1; month <= 12; month++ {
			period := MonthPeriod{
				MonthNumber: int32(month),
				Year:        year,
			}
			from, to := period.ConvertToTime()
			goal, _ := p.querier.GetSalesManagerGoalByGivenDateRangeAndSaleType(p.ctx, generated.GetSalesManagerGoalByGivenDateRangeAndSaleTypeParams{
				SalesManagerID: int32(smId),
				FromDate:       from,
				ToDate:         to,
				TypeID:         int32(saleType.Id),
			})

			sum, err := p.customQuerier.GetSalesManagerYearStatistic(p.ctx, GetSalesManagerYearStatisticParams{
				SalesManagerID: int32(smId),
				TypeId:         int32(saleType.Id),
				Year:           year,
				Month:          int32(month),
			})
			if err != nil {
				log.Println(err)
			}

			stat := MonthlyYearStatistic{
				SaleType: saleType,
				Month:    MonthNumber(month),
				Amount:   SaleAmount(sum.TotalAmount),
				Goal:     SaleAmount(goal),
			}
			result = append(result, stat)
		}
	}

	return &result, nil
}
