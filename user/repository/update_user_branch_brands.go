package repository

import (
	"context"
	"fmt"
	generated "zhasa2.0/db/sqlc"
)

type UpdateUserBrandsFunc func(userId int32, brands []int32) error

func NewUpdateUserBrandsFunc(ctx context.Context, store generated.UserStore) UpdateUserBrandsFunc {
	return func(userId int32, brands []int32) error {
		err := store.UpdateUserBrandsTX(ctx, userId, brands)
		if err != nil {
			fmt.Println(err)
		}
		return err
	}
}
