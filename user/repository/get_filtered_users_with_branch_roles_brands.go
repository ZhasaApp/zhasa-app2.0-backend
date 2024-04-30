package repository

import (
	"context"
	"strings"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/user/entities"
)

type GetFilteredUsersWithBranchBrands func(params generated.GetFilteredUsersWithBranchRolesBrandsParams) ([]entities.UserWithBranchBrands, int32, error)

func NewGetFilteredUsersWithBranchBrands(ctx context.Context, store generated.UserStore) GetFilteredUsersWithBranchBrands {
	return func(params generated.GetFilteredUsersWithBranchRolesBrandsParams) ([]entities.UserWithBranchBrands, int32, error) {
		rows, err := store.GetFilteredUsersWithBranchRolesBrands(ctx, params)
		users := make([]entities.UserWithBranchBrands, 0)
		if err != nil {
			return nil, 0, err
		}

		var total int32
		for _, row := range rows {
			users = append(users, entities.UserWithBranchBrands{
				Id:          row.ID,
				FirstName:   row.FirstName,
				LastName:    row.LastName,
				Phone:       row.Phone,
				BranchTitle: row.BranchTitle.String,
				Brands:      strings.Split(string(row.Brands), ", "),
				Role:        row.Role,
				Deleted:     row.Deleted,
			})
			total = int32(row.TotalCount)
		}
		return users, total, nil
	}
}
