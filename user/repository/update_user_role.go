package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type UpdateUserRoleFunc func(userId int32, role string) error

func NewUpdateUserRoleFunc(ctx context.Context, store generated.UserStore) UpdateUserRoleFunc {
	return func(userId int32, role string) error {
		params := generated.UpdateUserRoleParams{
			UserID:  userId,
			RoleKey: role,
		}
		return store.UpdateUserRole(ctx, params)
	}
}
