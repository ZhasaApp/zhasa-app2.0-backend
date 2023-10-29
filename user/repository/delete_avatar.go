package repository

import (
	"context"
	"fmt"
	generated "zhasa2.0/db/sqlc"
	. "zhasa2.0/user/entities"
)

type DeleteAvatarFunc func(userId UserId) error

func NewDeleteAvatarFunc(ctx context.Context, store generated.UserStore) DeleteAvatarFunc {
	return func(userId UserId) error {
		err := store.DeleteUserAvatar(ctx, int32(userId))
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	}
}
