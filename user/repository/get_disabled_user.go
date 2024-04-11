package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type CheckDisabledUserFunc func(userId int32) (bool, error)

func NewCheckDisabledUserFunc(ctx context.Context, store generated.UserStore) CheckDisabledUserFunc {
	return func(userId int32) (bool, error) {
		id, err := store.GetDisabledUser(ctx, userId)
		if err != nil {
			return false, err
		}
		return id != 0, err
	}
}
