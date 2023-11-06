package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type UpdateUserBranchFunc func(userId int32, branchId int32) error

func NewUpdateUserBranchFunc(ctx context.Context, store generated.UserStore) UpdateUserBranchFunc {
	return func(userId int32, branchId int32) error {
		params := generated.UpdateUserBranchParams{
			UserID:   userId,
			BranchID: branchId,
		}
		err := store.UpdateUserBranch(ctx, params)
		if err != nil {
			return err
		}
		return nil
	}
}
