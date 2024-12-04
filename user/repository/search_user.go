package repository

import (
	"context"
	"database/sql"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/user/entities"
)

type SearchUsersFunc func(search string) ([]*entities.User, error)

func NewSearchUsersFunc(ctx context.Context, store generated.UserStore) SearchUsersFunc {
	return func(search string) ([]*entities.User, error) {
		rows, err := store.SearchUsers(ctx, sql.NullString{String: search, Valid: true})
		if err != nil {
			return nil, err
		}

		users := make([]*entities.User, len(rows))
		for i, row := range rows {
			users[i] = &entities.User{
				Id:        row.ID,
				Phone:     entities.Phone(row.Phone),
				Avatar:    row.AvatarUrl,
				FirstName: row.FirstName,
				LastName:  row.LastName,
			}
		}

		return users, nil
	}
}
