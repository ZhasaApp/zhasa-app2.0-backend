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
	UploadUserAvatar(ctx context.Context, arg UploadUserAvatarParams) error
	DeleteUserAvatar(ctx context.Context, userID int32) error
	CreateUser(ctx context.Context, arg CreateUserParams) (int32, error)
	AddBrandToUser(ctx context.Context, arg AddBrandToUserParams) error
	AddRoleToUser(ctx context.Context, arg AddRoleToUserParams) error
	AddUserToBranch(ctx context.Context, arg AddUserToBranchParams) error
	CreateManagerTX(ctx context.Context, userId, branchId int32, brands []int32) error
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
