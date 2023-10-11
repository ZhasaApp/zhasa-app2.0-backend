package generated

import "context"

type UserStore interface {
	GetUserBranch(ctx context.Context, userID int32) (GetUserBranchRow, error)
	GetUsersOrderedByRatioForGivenBrand(ctx context.Context, arg GetUsersOrderedByRatioForGivenBrandParams) ([]GetUsersOrderedByRatioForGivenBrandRow, error)
}
