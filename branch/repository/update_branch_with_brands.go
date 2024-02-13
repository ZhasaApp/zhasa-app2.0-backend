package repository

import (
	"context"
	"zhasa2.0/branch/entities"
	generated "zhasa2.0/db/sqlc"
)

type UpdateBranchWithBrandsFunc func(branchWithBrands entities.BranchWithBrands) error

func NewUpdateBranchWithBrandsFunc(ctx context.Context, store generated.BranchStore) UpdateBranchWithBrandsFunc {
	return func(branchWithBrands entities.BranchWithBrands) error {
		return store.UpdateBranchWithBrandsTX(ctx, branchWithBrands)
	}
}
