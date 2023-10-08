package generated

import "context"

type UserBrandStore interface {
	GetUserBrandGoal(ctx context.Context, arg GetUserBrandGoalParams) (int64, error)
	GetUserBrand(ctx context.Context, arg GetUserBrandParams) (int32, error)
	InsertUserBrandRatio(ctx context.Context, arg InsertUserBrandRatioParams) error
	GetUserRank(ctx context.Context, arg GetUserRankParams) (int64, error)
}
