package repository

import (
	"context"
	"fmt"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/user/entities"
)

type GetUserByPhoneWithPasswordFunc func(phone entities.Phone) (*entities.AuthUser, error)

func NewGetUserByPhoneWithPasswordFunc(ctx context.Context, store generated.UserStore) GetUserByPhoneWithPasswordFunc {
	return func(phone entities.Phone) (*entities.AuthUser, error) {
		row, err := store.GetUserByPhoneWithPassword(ctx, phone.String())
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return &entities.AuthUser{
			Id:        row.ID,
			Phone:     entities.Phone(row.Phone),
			FirstName: row.FirstName,
			LastName:  row.LastName,
			Password:  row.Password.String,
		}, nil
	}
}
