package generated

import "context"

type UserStore interface {
	GetUserBranch(ctx context.Context, userID int32) (GetUserBranchRow, error)
	GetUsersOrderedByRatioForGivenBrand(ctx context.Context, arg GetUsersOrderedByRatioForGivenBrandParams) ([]GetUsersOrderedByRatioForGivenBrandRow, error)
	GetBranchUsersOrderedByRatioForGivenBrand(ctx context.Context, arg GetBranchUsersOrderedByRatioForGivenBrandParams) ([]GetBranchUsersOrderedByRatioForGivenBrandRow, error)
	GetBranchBrandUserByRole(ctx context.Context, arg GetBranchBrandUserByRoleParams) ([]GetBranchBrandUserByRoleRow, error)
	GetUsersByBranchBrandRole(ctx context.Context, arg GetUsersByBranchBrandRoleParams) ([]GetUsersByBranchBrandRoleRow, error)
}
