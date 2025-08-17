package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type AddUserTokenFunc func(userId int32, token string) error

func NewAddUserTokenFunc(ctx context.Context, store generated.UserStore) AddUserTokenFunc {
	return func(userId int32, token string) error {
		err := store.AddUserToken(ctx, generated.AddUserTokenParams{
			UserID: userId,
			Token:  token,
		})
		if err != nil {
			return err
		}
		return nil
	}
}
