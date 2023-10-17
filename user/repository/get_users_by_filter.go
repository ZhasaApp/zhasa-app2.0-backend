package repository

import (
	"context"
	"database/sql"
	"fmt"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/user/entities"
)

type GetUsersByBranchBrandRoleFunc func(branchId int32, brandId int32, roleId int32) ([]entities.User, error)

func NewGetUsersByBranchBrandRoleFunc(ctx context.Context, store generated.UserStore) GetUsersByBranchBrandRoleFunc {
	return func(branchId int32, brandId int32, roleId int32) ([]entities.User, error) {
		params := generated.GetUsersByBranchBrandRoleParams{
			BrandID:  brandId,
			BranchID: branchId,
			RoleID:   roleId,
		}
		rows, err := store.GetUsersByBranchBrandRole(ctx, params)
		users := make([]entities.User, 0)
		if err == sql.ErrNoRows {
			return users, nil
		}
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		for _, row := range rows {
			users = append(users, entities.User{
				Id:        row.ID,
				Phone:     "",
				Avatar:    row.AvatarUrl,
				FirstName: row.FirstName,
				LastName:  row.LastName,
				UserRole: entities.UserRole{
					Id:  roleId,
					Key: "",
				},
			})
		}
		return users, nil
	}
}
