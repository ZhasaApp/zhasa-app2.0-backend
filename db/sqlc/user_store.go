package generated

import "context"

type UserStore interface {
	GetUserBranch(ctx context.Context, userID int32) (GetUserBranchRow, error)
	GetUsersOrderedByRatioForGivenBrand(ctx context.Context, arg GetUsersOrderedByRatioForGivenBrandParams) ([]GetUsersOrderedByRatioForGivenBrandRow, error)
	GetBranchUsersOrderedByRatioForGivenBrand(ctx context.Context, arg GetBranchUsersOrderedByRatioForGivenBrandParams) ([]GetBranchUsersOrderedByRatioForGivenBrandRow, error)
	GetBranchBrandUserByRole(ctx context.Context, arg GetBranchBrandUserByRoleParams) ([]GetBranchBrandUserByRoleRow, error)
	GetUsersByBranchBrandRole(ctx context.Context, arg GetUsersByBranchBrandRoleParams) ([]GetUsersByBranchBrandRoleRow, error)
	GetUserByPhone(ctx context.Context, phone string) (UserAvatarView, error)
	GetUserById(ctx context.Context, id int32) (GetUserByIdRow, error)
	GetUsersWithoutRoles(ctx context.Context, search string) ([]GetUsersWithoutRolesRow, error)
	GetUsersWithBranchRolesBrands(ctx context.Context, arg GetUsersWithBranchRolesBrandsParams) ([]GetUsersWithBranchRolesBrandsRow, error)
	GetUsersWithBranchBrands(ctx context.Context, arg GetUsersWithBranchBrandsParams) ([]GetUsersWithBranchBrandsRow, error)
	UploadUserAvatar(ctx context.Context, arg UploadUserAvatarParams) error
	DeleteUserAvatar(ctx context.Context, userID int32) error
	CreateUser(ctx context.Context, arg CreateUserParams) (int32, error)
	CreateUserCode(ctx context.Context, arg CreateUserCodeParams) (int32, error)
	GetAuthCodeById(ctx context.Context, id int32) (UsersCode, error)
	AddBrandToUser(ctx context.Context, arg AddBrandToUserParams) error
	AddRoleToUser(ctx context.Context, arg AddRoleToUserParams) error
	AddUserToBranch(ctx context.Context, arg AddUserToBranchParams) error
	CreateManagerTX(ctx context.Context, userId, branchId int32, brands []int32) error
	UpdateUserBrandsTX(ctx context.Context, userId int32, brands []int32) error
	UpdateUser(ctx context.Context, arg UpdateUserParams) error
	UpdateUserBranch(ctx context.Context, params UpdateUserBranchParams) error
	AddDisabledUser(ctx context.Context, userID int32) error
	GetFilteredUsersWithBranchRolesBrands(ctx context.Context, arg GetFilteredUsersWithBranchRolesBrandsParams) ([]GetFilteredUsersWithBranchRolesBrandsRow, error)
	AddUserRole(ctx context.Context, arg AddUserRoleParams) error
	AddUserBranch(ctx context.Context, arg AddUserBranchParams) error
}

func (db *DBStore) CreateManagerTX(ctx context.Context, userId, branchId int32, brands []int32) error {
	return db.execTx(ctx, func(queries *Queries) error {
		const managerRoleId = 2
		err := queries.AddRoleToUser(ctx, AddRoleToUserParams{
			UserID: userId,
			RoleID: managerRoleId,
		})
		if err != nil {
			return err
		}
		err = queries.AddUserToBranch(ctx, AddUserToBranchParams{
			UserID:   userId,
			BranchID: branchId,
		})
		if err != nil {
			return err
		}
		for _, brandId := range brands {
			err = queries.AddBrandToUser(ctx, AddBrandToUserParams{
				UserID:  userId,
				BrandID: brandId,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *DBStore) UpdateUserBrandsTX(ctx context.Context, userId int32, brands []int32) error {
	return db.execTx(ctx, func(queries *Queries) error {
		err := queries.DeleteUserBrandByUserId(ctx, userId)
		if err != nil {
			return err
		}

		for _, brandId := range brands {
			err = queries.AddBrandToUser(ctx, AddBrandToUserParams{
				UserID:  userId,
				BrandID: brandId,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}
