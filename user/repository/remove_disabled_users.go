package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type RemoveDisabledUsersFunc func(userIds []int32) error

func NewRemoveDisabledUsersFunc(ctx context.Context, store generated.UserStore) RemoveDisabledUsersFunc {
	return func(userIds []int32) error {
		err := store.DeleteDisabledUsers(ctx, userIds)
		if err != nil {
			return err
		}
		return nil
	}
}
