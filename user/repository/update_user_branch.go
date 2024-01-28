package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type UpdateUserBranchFunc func(userId int32, branchId *int32) error

func NewUpdateUserBranchFunc(ctx context.Context, store generated.UserStore) UpdateUserBranchFunc {
	return func(userId int32, branchId *int32) error {
		if branchId == nil {
			return store.DeleteUserBranchByUserId(ctx, userId)
		}
		err := store.UpdateUserBranchTX(ctx, userId, *branchId)
		if err != nil {
			return err
		}
		return nil
	}
}
