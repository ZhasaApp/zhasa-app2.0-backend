package repository

import (
	"context"
	"errors"
	"fmt"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/user/entities"
)

type CreateUserFunc func(firstName, lastName string, phone entities.Phone) (int32, error)

func NewCreateUserFunc(ctx context.Context, store generated.UserStore) CreateUserFunc {
	return func(firstName, lastName string, phone entities.Phone) (int32, error) {
		id, err := store.CreateUser(ctx, generated.CreateUserParams{
			FirstName: firstName,
			LastName:  lastName,
			Phone:     phone.String(),
		})

		if err != nil {
			fmt.Println(err)
			return 0, errors.New("error creating user")
		}
		return id, nil
	}
}
