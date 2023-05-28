package repository

import (
	"context"
	"database/sql"
	"log"
	"time"
	. "zhasa2.0/base"
	. "zhasa2.0/branch/entities"
	. "zhasa2.0/branch/repository"
	. "zhasa2.0/db/hand-made"
	. "zhasa2.0/manager/entities"
	. "zhasa2.0/statistic"
	. "zhasa2.0/user/entities"
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
	params := GetRankedSalesManagersParams{
		FromDate: from,
		ToDate:   to,
		Limit:    pagination.PageSize,
		Offset:   pagination.Page,
	}

	data, err := j.cQ.GetRankedSalesManagers(j.ctx, params)

	log.Println(params)

	result := make([]SalesManager, 0)
	if err == sql.ErrNoRows {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	for _, row := range data {
		branch, err := j.GetBranchById(BranchId(row.BranchID))

		if err != nil {
			return nil, err
		}

		result = append(result, SalesManager{
			Id:          SalesManagerId(row.SalesManagerID),
			UserId:      UserId(row.UserId),
			FirstName:   row.FirstName,
			LastName:    row.LastName,
			AvatarUrl:   "",
			Branch:      *branch,
			Ratio:       Percent(row.Ratio),
			RatingPlace: RatingPlace(row.RatingPosition),
		})
	}
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
