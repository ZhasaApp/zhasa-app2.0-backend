package repository

import (
	"context"
	"database/sql"
	"zhasa2.0/base"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/user/entities"
)

type GetUsersByRoleFunc func(search, key string, pagination base.Pagination) ([]entities.UserWithBrands, int32, error)

func NewGetUsersByRoleFunc(ctx context.Context, store generated.UserStore) GetUsersByRoleFunc {
	return func(search, key string, pagination base.Pagination) ([]entities.UserWithBrands, int32, error) {
		params := generated.GetUsersWithBranchRolesBrandsParams{
			Search: search,
			Key:    key,
			Limit:  pagination.PageSize,
			Offset: pagination.GetOffset(),
		}
		rows, err := store.GetUsersWithBranchRolesBrands(ctx, params)
		users := make([]entities.UserWithBrands, 0)
		if err == sql.ErrNoRows {
			return users, 0, nil
		}
		if err != nil {
			return nil, 0, err
		}

		var total int32
		for _, row := range rows {
			users = append(users, entities.UserWithBrands{
				Id:          row.ID,
				FirstName:   row.FirstName,
				LastName:    row.LastName,
				Phone:       row.Phone,
				BranchTitle: row.BranchTitle,
				Brands:      string(row.Brands),
				IsActive:    row.IsActive,
			})
			total = int32(row.TotalCount)
		}
		return users, total, nil
	}
}
