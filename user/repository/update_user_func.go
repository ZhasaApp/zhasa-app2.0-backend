package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/user/entities"
)

type UpdateUserFunc func(userId int32, firstName, lastName entities.Name, phone entities.Phone) error

func NewUpdateUserFunc(ctx context.Context, store generated.UserStore) UpdateUserFunc {
	return func(userId int32, firstName, lastName entities.Name, phone entities.Phone) error {
		return store.UpdateUser(ctx, generated.UpdateUserParams{
			FirstName: firstName.String(),
			LastName:  lastName.String(),
			Phone:     phone.String(),
			ID:        userId,
		})
	}
}
