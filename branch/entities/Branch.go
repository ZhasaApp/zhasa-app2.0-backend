package entities

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type BranchId int32

type Branch struct {
	BranchId    BranchId
	Title       string
	Description string
}

type BranchesMap map[BranchId]*Branch

type BranchRepository interface {
	GetBranch(id BranchId) (*Branch, error)
}

type DBBranchRepository struct {
	ctx     context.Context
	querier generated.Querier
	cache   BranchesMap
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
		Title:       branchDb.Title,
		Description: branchDb.Description.String,
	}

	br.cache[id] = newBranch
	return newBranch, nil
}
