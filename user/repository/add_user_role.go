package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type AddUserRoleFunc func(userId int32, roleKey string) error

func NewAddUserRoleFunc(ctx context.Context, store generated.UserStore) AddUserRoleFunc {
	return func(userId int32, roleKey string) error {
		params := generated.AddUserRoleParams{
			UserID:  userId,
			RoleKey: roleKey,
		}
		return store.AddUserRole(ctx, params)
	}
}
