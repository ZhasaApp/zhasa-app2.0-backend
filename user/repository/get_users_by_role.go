package repository

import (
	"context"
	"database/sql"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/user/entities"
)

type GetUsersByRoleFunc func(key string) ([]entities.UserWithBrands, error)

func NewGetUsersByRoleFunc(ctx context.Context, store generated.UserStore) GetUsersByRoleFunc {
	return func(key string) ([]entities.UserWithBrands, error) {
		rows, err := store.GetUsersWithBranchRolesBrands(ctx, key)
		users := make([]entities.UserWithBrands, 0)
		if err == sql.ErrNoRows {
			return users, nil
		}
		if err != nil {
			return nil, err
		}

		for _, row := range rows {
			users = append(users, entities.UserWithBrands{
				Id:          row.ID,
				FirstName:   row.FirstName,
				LastName:    row.LastName,
				BranchTitle: row.BranchTitle,
				Brands:      string(row.Brands),
			})
		}
		return users, nil
	}
}
