package repository

import (
	"context"
	"database/sql"
	"fmt"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/user/entities"
)

type GetUsersWithoutRolesFunc func(search string) ([]entities.BaseUser, error)

func NewGetUsersWithoutRolesFunc(ctx context.Context, store generated.UserStore) GetUsersWithoutRolesFunc {
	return func(search string) ([]entities.BaseUser, error) {
		rows, err := store.GetUsersWithoutRoles(ctx, search)
		users := make([]entities.BaseUser, 0)
		if err == sql.ErrNoRows {
			return users, nil
		}
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		for _, row := range rows {
			users = append(users, entities.BaseUser{
				Id:        row.ID,
				FirstName: row.FirstName,
				LastName:  row.LastName,
			})
		}
		return users, nil
	}
}
