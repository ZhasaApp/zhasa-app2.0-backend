package generated

import "context"

type BranchBrandStore interface {
	GetBranchBrands(ctx context.Context, branchID int32) ([]GetBranchBrandsRow, error)
}
