package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
	. "zhasa2.0/base"
	. "zhasa2.0/branch/entities"
	handmade "zhasa2.0/db/hand-made"
	generated "zhasa2.0/db/sqlc"
	. "zhasa2.0/manager/entities"
	. "zhasa2.0/sale/entities"
	. "zhasa2.0/sale/repository"
	. "zhasa2.0/statistic"
	. "zhasa2.0/statistic/entities"
	. "zhasa2.0/user/entities"
)

type BranchRepository interface {
	CreateBranch(request CreateBranchRequest) error
	GetBranchById(id BranchId) (*Branch, error)
	GetBranches() ([]Branch, error)
	GetBranchYearMonthlyStatistic(bId BranchId, year int32) (*[]MonthlyYearStatistic, error)
	GetBranchSalesSums(from, to time.Time, branchId BranchId) (*SaleSumByType, error)
	GetBranchGoal(from, to time.Time, branchId BranchId, typeId SaleTypeId) (SaleAmount, error)
	GetBranchRankedSalesManagers(from, to time.Time, branchId BranchId, pagination Pagination) (*[]SalesManager, error)
}

type DBBranchRepository struct {
	ctx           context.Context
	querier       generated.Querier
	customQuerier handmade.CustomQuerier
	SaleTypeRepository
	cache BranchesMap
}

func (br DBBranchRepository) GetBranchRankedSalesManagers(from, to time.Time, branchId BranchId, pagination Pagination) (*[]SalesManager, error) {
	params := generated.GetOrderedSalesManagersOfBranchParams{
		FromDate: from,
		ToDate:   to,
		BranchID: int32(branchId),
		Limit:    pagination.PageSize,
		Offset:   pagination.Page,
	}

	branch, err := br.GetBranchById(branchId)

	if err != nil {
		return nil, err
	}

	data, err := br.querier.GetOrderedSalesManagersOfBranch(br.ctx, params)

	result := make([]SalesManager, 0)
	if err == sql.ErrNoRows {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	for index, row := range data {
		log.Println(row)
		result = append(result, SalesManager{
			Id:          SalesManagerId(row.SalesManagerID),
			UserId:      UserId(row.UserID),
			FirstName:   row.FirstName,
			LastName:    row.LastName,
			AvatarUrl:   row.AvatarUrl,
			Branch:      *branch,
			Ratio:       Percent(row.Ratio).GetRounded(),
			RatingPlace: RatingPlace((pagination.Page)*pagination.PageSize + int32(index) + int32(1)),
		})
	}
	return &result, nil
}

func (br DBBranchRepository) GetBranchSalesSums(from, to time.Time, branchId BranchId) (*SaleSumByType, error) {
	types, err := br.GetSaleTypes()
	if err != nil {
		return nil, err
	}

	sums := make([]SumsByTypeRow, 0)

	for _, t := range *types {
		arg := handmade.GetBranchSumByTypeParams{
			SaleDate:   from,
			SaleDate_2: to,
			ID:         int32(t.Id),
			BranchID:   int32(branchId),
		}

		res, err := br.customQuerier.GetBranchSumByType(br.ctx, arg)
		if err != nil {
			return nil, err
		}

		sums = append(sums, SumsByTypeRow{
			SaleTypeID:    int32(t.Id),
			SaleTypeTitle: t.Title,
			TotalSales:    res.TotalSales,
		})
	}

	result := br.MapSalesSumsByType(sums)
	return &result, err
}

func (br DBBranchRepository) GetBranchGoal(from, to time.Time, branchId BranchId, typeId SaleTypeId) (SaleAmount, error) {
	arg := generated.GetBranchGoalByGivenDateRangeParams{
		BranchID: int32(branchId),
		FromDate: from,
		ToDate:   to,
		TypeID:   int32(typeId),
	}
	log.Println(arg)
	data, err := br.querier.GetBranchGoalByGivenDateRange(br.ctx, arg)
	if err == sql.ErrNoRows {
		log.Println(err)
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return SaleAmount(data), nil
}

func (br DBBranchRepository) GetBranchYearMonthlyStatistic(bId BranchId, year int32) (*[]MonthlyYearStatistic, error) {

	result := make([]MonthlyYearStatistic, 0)

	saleTypes, err := br.GetSaleTypes()
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
			goal, _ := br.querier.GetBranchGoalByGivenDateRange(br.ctx, generated.GetBranchGoalByGivenDateRangeParams{
				BranchID: int32(bId),
				FromDate: from,
				ToDate:   to,
				TypeID:   int32(saleType.Id),
			})

			sum, err := br.customQuerier.GetBranchYearStatistic(br.ctx, handmade.GetBranchYearStatisticParams{
				BranchId: int32(bId),
				TypeId:   int32(saleType.Id),
				Year:     year,
				Month:    int32(month),
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

func NewBranchRepository(ctx context.Context, querier generated.Querier, customQ handmade.CustomQuerier, sTypeRepo SaleTypeRepository) BranchRepository {
	cache := make(BranchesMap)
	return DBBranchRepository{
		ctx:                ctx,
		querier:            querier,
		customQuerier:      customQ,
		SaleTypeRepository: sTypeRepo,
		cache:              cache,
	}
}

func (br DBBranchRepository) CreateBranch(request CreateBranchRequest) error {
	params := generated.CreateBranchParams{
		Title:       string(request.Title),
		Description: string(request.Description),
		BranchKey:   string(request.Key),
	}
	return br.querier.CreateBranch(br.ctx, params)
}

func (br DBBranchRepository) GetBranchById(id BranchId) (*Branch, error) {
	branch, found := br.cache[id]

	if found {
		return branch, nil
	}

	branchDb, err := br.querier.GetBranchById(br.ctx, int32(id))
	if err != nil {
		return nil, errors.New("no branch found for given id")
	}

	newBranch := &Branch{
		BranchId:    BranchId(branchDb.ID),
		Title:       BranchTitle(branchDb.Title),
		Description: BranchDescription(branchDb.Description),
		Key:         BranchKey(branchDb.BranchKey),
	}

	br.cache[id] = newBranch
	return newBranch, nil
}

func (br DBBranchRepository) GetBranches() ([]Branch, error) {

	branchList := make([]Branch, 0)

	return branchList, nil
}
