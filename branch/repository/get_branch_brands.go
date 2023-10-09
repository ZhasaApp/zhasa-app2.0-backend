package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type GetBranchBrandsFunc func(branchId int32) ([]generated.GetBranchBrandsRow, error)

func NewGetBranchBrandsFunc(ctx context.Context, store generated.BranchBrandStore) GetBranchBrandsFunc {
	return func(branchId int32) ([]generated.GetBranchBrandsRow, error) {
		brands, err := store.GetBranchBrands(ctx, branchId)
		if err != nil {
			return nil, err
		}
		return brands, nil
	}
}
