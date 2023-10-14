package generated

import "context"

type BranchStore interface {
	GetBranchBrand(ctx context.Context, arg GetBranchBrandParams) (int32, error)
	GetBranchBrandGoalByGivenDateRange(ctx context.Context, arg GetBranchBrandGoalByGivenDateRangeParams) (int64, error)
	GetBranchById(ctx context.Context, id int32) (Branch, error)
	GetBranchesByBrandId(ctx context.Context, brandID int32) ([]GetBranchesByBrandIdRow, error)
	SetBranchBrandGoal(ctx context.Context, arg SetBranchBrandGoalParams) error
}
