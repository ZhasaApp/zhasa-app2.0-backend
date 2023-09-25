package repository

import (
	"context"
	"errors"
	"time"
	. "zhasa2.0/branch/entities"
	handmade "zhasa2.0/db/hand-made"
	generated "zhasa2.0/db/sqlc"
	. "zhasa2.0/sale/entities"
	. "zhasa2.0/sale/repository"
	. "zhasa2.0/statistic"
	. "zhasa2.0/statistic/entities"
)

type BranchRepository interface {
	CreateBranch(request CreateBranchRequest) error
	GetBranchById(id BranchId) (*Branch, error)
	GetBranches(period Period) ([]Branch, error)
	GetBranchYearMonthlyStatistic(bId BranchId, year int32) (*[]MonthlyYearStatistic, error)
	GetBranchSalesSums(from, to time.Time, branchId BranchId) (*SaleSumByType, error)
	GetBranchGoal(from, to time.Time, branchId BranchId, typeId SaleTypeId) (SaleAmount, error)
}

type DBBranchRepository struct {
	ctx           context.Context
	querier       generated.Querier
	customQuerier handmade.CustomQuerier
	SaleTypeRepository
	cache BranchesMap
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
	//arg := generated.GetBranchGoalByGivenDateRangeParams{
	//	//	BranchID: int32(branchId),
	//	FromDate: from,
	//	ToDate:   to,
	//	//	TypeID:   int32(typeId),
	//}
	//data, err := br.querier.GetBranchGoalByGivenDateRange(br.ctx, arg)
	//if err == sql.ErrNoRows {
	//	log.Println(err)
	//	return 0, nil
	//}
	//if err != nil {
	//	return 0, err
	//}
	//
	//return SaleAmount(data), nil
	return 0, nil
}

func (br DBBranchRepository) GetBranchYearMonthlyStatistic(bId BranchId, year int32) (*[]MonthlyYearStatistic, error) {

	result := make([]MonthlyYearStatistic, 0)

	//saleTypes, err := br.GetSaleTypes()
	//if err != nil {
	//	return nil, errors.New("error while getting sale types")
	//}
	//
	//for _, saleType := range *saleTypes {
	//	for month := 1; month <= 12; month++ {
	//		period := MonthPeriod{
	//			MonthNumber: int32(month),
	//			Year:        year,
	//		}
	//		from, to := period.ConvertToTime()
	//		goal, _ := br.querier.GetBranchGoalByGivenDateRange(br.ctx, generated.GetBranchGoalByGivenDateRangeParams{
	//			//	BranchID: int32(bId),
	//			FromDate: from,
	//			ToDate:   to,
	//			//		TypeID:   int32(saleType.Id),
	//		})
	//
	//		sum, err := br.customQuerier.GetBranchYearStatistic(br.ctx, handmade.GetBranchYearStatisticParams{
	//			BranchId: int32(bId),
	//			TypeId:   int32(saleType.Id),
	//			Year:     year,
	//			Month:    int32(month),
	//		})
	//
	//		if err != nil {
	//			log.Println(err)
	//		}
	//
	//		stat := MonthlyYearStatistic{
	//			SaleType: saleType,
	//			Month:    MonthNumber(month),
	//			Amount:   SaleAmount(sum.TotalAmount),
	//			Goal:     SaleAmount(goal),
	//		}
	//		result = append(result, stat)
	//	}
	//}

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
	}

	br.cache[id] = newBranch
	return newBranch, nil
}

func (br DBBranchRepository) GetBranches(period Period) ([]Branch, error) {

	branchList := make([]Branch, 0)

	//rows, err := br.querier.GetBranches(br.ctx)
	//
	//if err == sql.ErrNoRows {
	//	return branchList, nil
	//}
	//
	//if err != nil {
	//	log.Println(err)
	//	return nil, err
	//}
	//
	//saleTypes, err := br.GetSaleTypes()
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//from, to := period.ConvertToTime()
	//
	//for _, row := range rows {
	//	ratioRows := make([]RatioRow, 0)
	//
	//	for _, sType := range *saleTypes {
	//		goal, _ := br.querier.GetBranchGoalByGivenDateRange(br.ctx, generated.GetBranchGoalByGivenDateRangeParams{
	//			BranchID: row.ID,
	//			FromDate: from,
	//			ToDate:   to,
	//			TypeID:   int32(sType.Id),
	//		})
	//
	//		sum, _ := br.customQuerier.GetBranchSumByType(br.ctx, handmade.GetBranchSumByTypeParams{
	//			BranchID:   row.ID,
	//			SaleDate:   from,
	//			SaleDate_2: to,
	//			ID:         int32(sType.Id),
	//		})
	//
	//		ratioRows = append(ratioRows, RatioRow{
	//			Amount:  SaleAmount(sum.TotalSales),
	//			Goal:    SaleAmount(goal),
	//			Gravity: sType.Gravity,
	//		})
	//	}
	//	percent := CalculateRatio(ratioRows)
	//	fmt.Println(row.Title+": ", percent)
	//	branchList = append(branchList, Branch{
	//		BranchId:        BranchId(row.ID),
	//		Title:           BranchTitle(row.Title),
	//		Description:     BranchDescription(row.Description),
	//		Key:             BranchKey(row.BranchKey),
	//		GoalAchievement: Percent(percent).GetRounded(),
	//	})
	//}

	return branchList, nil
}
