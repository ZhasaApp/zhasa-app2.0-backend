package generated

import "context"

type BrandStore interface {
	GetBranchBrands(ctx context.Context, branchID int32) ([]GetBranchBrandsRow, error)
	GetBrands(ctx context.Context, arg GetBrandsParams) ([]GetBrandsRow, error)
	GetUserBrands(ctx context.Context, userID int32) ([]GetUserBrandsRow, error)
}
