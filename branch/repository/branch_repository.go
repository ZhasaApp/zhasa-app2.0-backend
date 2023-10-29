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
)

type BranchRepository interface {
	CreateBranch(request CreateBranchRequest) error
	GetBranchById(id BranchId) (*Branch, error)
	GetBranchSalesSums(from, to time.Time, branchId BranchId) (*SaleSumByType, error)
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
			ID:         t.Id,
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
		BranchId:    branchDb.ID,
		Title:       branchDb.Title,
		Description: branchDb.Description,
	}

	br.cache[id] = newBranch
	return newBranch, nil
}
