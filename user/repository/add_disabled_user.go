package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type AddDisabledUserFunc func(userId int32) error

func NewAddDisabledUserFunc(ctx context.Context, store generated.UserStore) AddDisabledUserFunc {
	return func(userId int32) error {
		err := store.AddDisabledUser(ctx, userId)
		if err != nil {
			return err
		}
		return nil
	}
}
