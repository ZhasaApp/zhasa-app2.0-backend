package repository

import (
	"context"
	"database/sql"
	generated "zhasa2.0/db/sqlc"
)

type UpdateUserProfileAbout func(id int32, about *string) error

func NewUpdateUserProfileAbout(ctx context.Context, store generated.UserStore) UpdateUserProfileAbout {
	return func(id int32, about *string) error {
		aboutString := sql.NullString{Valid: false}
		if about != nil {
			aboutString = sql.NullString{Valid: true, String: *about}
		}
		err := store.UpdateUserAbout(ctx, generated.UpdateUserAboutParams{
			ID:    id,
			About: aboutString,
		})
		if err != nil {
			return err
		}
		return nil
	}
}
