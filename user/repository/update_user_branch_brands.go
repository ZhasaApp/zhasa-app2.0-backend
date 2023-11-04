package repository

import (
	"context"
	"fmt"
	generated "zhasa2.0/db/sqlc"
)

type UpdateUserBranchBrandsFunc func(userId, branchId int32, brands []int32) error

func NewUpdateUserBranchBrandsFunc(ctx context.Context, store generated.UserStore) UpdateUserBranchBrandsFunc {
	return func(userId, branchId int32, brands []int32) error {
		err := store.UpdateUserBranchBrandsTX(ctx, userId, branchId, brands)
		if err != nil {
			fmt.Println(err)
		}
		return err
	}
}
