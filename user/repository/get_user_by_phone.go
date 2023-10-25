package repository

import (
	"context"
	"fmt"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/user/entities"
)

type GetUserByPhoneFunc func(phone entities.Phone) (*entities.User, error)

func NewGetUserByPhoneFunc(ctx context.Context, store generated.UserStore) GetUserByPhoneFunc {
	return func(phone entities.Phone) (*entities.User, error) {
		row, err := store.GetUserByPhone(ctx, phone.String())
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return &entities.User{
			Id:        row.ID,
			Phone:     entities.Phone(row.Phone),
			Avatar:    row.AvatarUrl,
			FirstName: row.FirstName,
			LastName:  row.LastName,
			UserRole: entities.UserRole{
				Id:  row.RoleID,
				Key: row.RoleKey,
			},
		}, nil
	}
}
