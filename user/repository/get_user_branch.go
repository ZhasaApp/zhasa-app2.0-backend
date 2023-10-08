package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type GetUserBranchFunc func(userId int32) (generated.GetUserBranchRow, error)

func NewGetUserBranchFunc(ctx context.Context, store generated.UserStore) GetUserBranchFunc {
	return func(userId int32) (generated.GetUserBranchRow, error) {
		return store.GetUserBranch(ctx, userId)
	}
}
