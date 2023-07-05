package repository

import (
	"context"
	"time"
	. "zhasa2.0/base"
	. "zhasa2.0/branch/entities"
	. "zhasa2.0/branch/repository"
	. "zhasa2.0/db/hand-made"
	. "zhasa2.0/manager/entities"
	. "zhasa2.0/statistic"
)

type StatisticRepository interface {
	ProvideSums(salesManagerId SalesManagerId, from time.Time, to time.Time) (*SaleSumByType, error)
}

type RankingsRepository interface {
	ProvideRankedManagers(pagination Pagination, from time.Time, to time.Time) (*[]SalesManager, error)
	ProvideRankedBranches(pagination Pagination) (*[]Branch, error)
}

type DBRankingsRepository struct {
	ctx context.Context
	cQ  CustomQuerier
	BranchRepository
}

func (j DBRankingsRepository) ProvideRankedManagers(pagination Pagination, from time.Time, to time.Time) (*[]SalesManager, error) {
	result := make([]SalesManager, 0)

	return &result, nil
}

func (j DBRankingsRepository) ProvideRankedBranches(pagination Pagination) (*[]Branch, error) {
	branches := make([]Branch, 0)
	return &branches, nil
}

func NewRankingsRepository(ctx context.Context, cQ CustomQuerier, branchRepo BranchRepository) RankingsRepository {
	return DBRankingsRepository{
		ctx:              ctx,
		cQ:               cQ,
		BranchRepository: branchRepo,
	}
}
