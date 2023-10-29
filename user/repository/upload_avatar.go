package repository

import (
	"context"
	"fmt"
	generated "zhasa2.0/db/sqlc"
	. "zhasa2.0/user/entities"
)

type UploadAvatarFunc func(userId UserId, avatarUrl string) error

func NewUploadAvatarFunc(ctx context.Context, store generated.UserStore) UploadAvatarFunc {
	return func(userId UserId, avatarUrl string) error {
		err := store.UploadUserAvatar(ctx, generated.UploadUserAvatarParams{
			UserID:    int32(userId),
			AvatarUrl: avatarUrl,
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	}
}
