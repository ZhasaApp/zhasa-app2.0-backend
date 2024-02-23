package repository

import (
	"context"
	"zhasa2.0/branch/entities"
	generated "zhasa2.0/db/sqlc"
)

type CreateBranchWithBrandsFunc func(branch entities.BranchWithBrands) error

func NewCreateBranchWithBrandsFunc(ctx context.Context, store generated.BranchStore) CreateBranchWithBrandsFunc {
	return func(branch entities.BranchWithBrands) error {
		return store.AddBranchWithBrandsTX(ctx, branch)
	}
}
