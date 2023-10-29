package repository

import (
	"context"
	"fmt"
	generated "zhasa2.0/db/sqlc"
	. "zhasa2.0/user/entities"
)

type GetUserByIdFunc func(userId int32) (*User, error)

func NewGetUserByIdFunc(ctx context.Context, store generated.UserStore) GetUserByIdFunc {
	return func(userId int32) (*User, error) {
		row, err := store.GetUserById(ctx, userId)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return &User{
			Id:        row.ID,
			Phone:     Phone(row.Phone),
			Avatar:    row.AvatarUrl,
			FirstName: row.FirstName,
			LastName:  row.LastName,
			UserRole: UserRole{
				Id:  row.RoleID,
				Key: row.Key,
			},
		}, err
	}
}
