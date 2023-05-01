package repository

import (
	"context"
	"zhasa2.0/branch/entities"
	generated "zhasa2.0/db/sqlc"
)

type BranchRepository interface {
	CreateBranch(request entities.CreateBranchRequest) error
	GetBranch(id entities.BranchId) (*entities.Branch, error)
}

type DBBranchRepository struct {
	ctx     context.Context
	querier generated.Querier
	cache   entities.BranchesMap
}

func NewBranchRepository(ctx context.Context, querier generated.Querier) BranchRepository {
	cache := make(entities.BranchesMap)
	return DBBranchRepository{
		ctx:     ctx,
		querier: querier,
		cache:   cache,
	}
}

func (br DBBranchRepository) CreateBranch(request entities.CreateBranchRequest) error {
	params := generated.CreateBranchParams{
		Title:       string(request.Title),
		Description: string(request.Description),
	}
	return br.querier.CreateBranch(br.ctx, params)
}

func (br DBBranchRepository) GetBranch(id entities.BranchId) (*entities.Branch, error) {
	branch, found := br.cache[id]

	if found {
		return branch, nil
	}

	branchDb, err := br.querier.GetBranchById(br.ctx, int32(id))
	if err != nil {
		return nil, err
	}

	newBranch := &entities.Branch{
		BranchId:    entities.BranchId(branchDb.ID),
		Title:       entities.BranchTitle(branchDb.Title),
		Description: entities.BranchDescription(branchDb.Description),
	}

	br.cache[id] = newBranch
	return newBranch, nil
}
