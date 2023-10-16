package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/user/entities"
)

type GetUserByBranchBrandRoleFunc func(branchId int32, brandId int32, roleId int32) ([]entities.User, error)

func NewGetUserByBranchBrandRoleFunc(ctx context.Context, store generated.UserStore) GetUserByBranchBrandRoleFunc {
	return func(branchId int32, brandId int32, roleId int32) ([]entities.User, error) {
		rows, err := store.GetBranchBrandUserByRole(ctx, generated.GetBranchBrandUserByRoleParams{
			BranchID: branchId,
			BrandID:  brandId,
			RoleID:   roleId,
		})
		if err != nil {
			return nil, err
		}

		users := make([]entities.User, 0)
		for _, row := range rows {
			user := entities.User{
				Id:        row.ID,
				Avatar:    row.AvatarUrl,
				FirstName: row.FirstName,
				LastName:  row.LastName,
				UserRole: entities.UserRole{
					Id: roleId,
				},
			}
			users = append(users, user)
		}

		return users, nil
	}
}
