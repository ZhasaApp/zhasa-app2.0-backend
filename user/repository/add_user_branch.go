package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type AddUserBranchFunc func(userId int32, branchId int32) error

func NewAddUserBranchFunc(ctx context.Context, store generated.UserStore) AddUserBranchFunc {
	return func(userId int32, branchId int32) error {
		params := generated.AddUserBranchParams{
			UserID:   userId,
			BranchID: branchId,
		}
		return store.AddUserBranch(ctx, params)
	}
}
