package generated

import "context"

type BranchStore interface {
	GetBranchBrand(ctx context.Context, arg GetBranchBrandParams) (int32, error)
	GetBranchBrandGoalByGivenDateRange(ctx context.Context, arg GetBranchBrandGoalByGivenDateRangeParams) (int64, error)
}