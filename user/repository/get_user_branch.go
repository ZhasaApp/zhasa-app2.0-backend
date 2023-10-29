package repository

import (
	"context"
	"fmt"
	generated "zhasa2.0/db/sqlc"
)

type GetUserBranchFunc func(userId int32) (*generated.GetUserBranchRow, error)

func NewGetUserBranchFunc(ctx context.Context, store generated.UserStore) GetUserBranchFunc {
	return func(userId int32) (*generated.GetUserBranchRow, error) {
		branch, err := store.GetUserBranch(ctx, userId)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return &branch, nil
	}
}
