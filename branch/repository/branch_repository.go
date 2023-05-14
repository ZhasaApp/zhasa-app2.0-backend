package repository

import (
	"context"
	"database/sql"
	"time"
	. "zhasa2.0/branch/entities"
	handmade "zhasa2.0/db/hand-made"
	generated "zhasa2.0/db/sqlc"
	. "zhasa2.0/sale/entities"
	. "zhasa2.0/sale/repository"
	. "zhasa2.0/statistic/entities"
)

type BranchRepository interface {
	CreateBranch(request CreateBranchRequest) error
	GetBranch(id BranchId) (*Branch, error)
	GetBranches() ([]Branch, error)
	GetBranchYearMonthlyStatistic(id BranchId, year int32) (*[]MonthlyYearStatistic, error)
}

type DBBranchRepository struct {
	ctx           context.Context
	querier       generated.Querier
	customQuerier handmade.CustomQuerier
	SaleTypeRepository
	cache BranchesMap
}

func (br DBBranchRepository) GetBranchYearMonthlyStatistic(id BranchId, year int32) (*[]MonthlyYearStatistic, error) {
	params := handmade.GetBranchYearStatisticParams{
		BranchID:   int32(id),
		YearNumber: year,
	}
	data, err := br.customQuerier.GetBranchYearStatistic(br.ctx, params)

	ans := make([]MonthlyYearStatistic, 0)
	if err == sql.ErrNoRows {
		return &ans, nil
	}
	if err != nil {
		return nil, err
	}

	for _, item := range data {
		saleType, err := br.GetSaleType(SaleTypeId(item.SaleTypeId))
		if err != nil {
			return nil, err
		}
		ans = append(ans, MonthlyYearStatistic{
			SaleType: *saleType,
			Month:    MonthNumber(item.MonthNumber),
			Amount:   SaleAmount(item.TotalAmount),
		})
	}
	return &ans, nil
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

func (br DBBranchRepository) GetBranch(id BranchId) (*Branch, error) {
	branch, found := br.cache[id]

	if found {
		return branch, nil
	}

	branchDb, err := br.querier.GetBranchById(br.ctx, int32(id))
	if err != nil {
		return nil, err
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
	params := generated.GetBranchesByRatingParams{
		SaleDate:   time.Now(),
		SaleDate_2: time.Now().Add(time.Hour),
	}
	branches, err := br.querier.GetBranchesByRating(br.ctx, params)

	if err != nil {
		return nil, err
	}

	branchList := make([]Branch, 0)

	for _, br := range branches {
		branch := Branch{
			BranchId:    BranchId(br.BranchID),
			Title:       NewBranchTitle(br.BranchTitle),
			Description: NewBranchDescription(br.Description),
			Key:         NewBranchKey(br.BranchKey),
		}
		branchList = append(branchList, branch)
	}
	return branchList, err
}
