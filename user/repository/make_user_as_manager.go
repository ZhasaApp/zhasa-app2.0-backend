package repository

import (
	"context"
	"fmt"
	generated "zhasa2.0/db/sqlc"
)

type MakeUserAsManagerFunc func(userId, branchId int32, brands []int32) error

func NewMakeUserAsManagerFunc(ctx context.Context, store generated.UserStore) MakeUserAsManagerFunc {
	return func(userId, branchId int32, brands []int32) error {
		err := store.CreateManagerTX(ctx, userId, branchId, brands)
		if err != nil {
			fmt.Println(err)
		}
		return err
	}
}
