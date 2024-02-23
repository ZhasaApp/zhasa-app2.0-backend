package generated

import (
	"context"
	"zhasa2.0/branch/entities"
)

type BranchStore interface {
	GetBranchBrand(ctx context.Context, arg GetBranchBrandParams) (int32, error)
	GetBranchBrandGoalByGivenDateRange(ctx context.Context, arg GetBranchBrandGoalByGivenDateRangeParams) (int64, error)
	GetBranchById(ctx context.Context, id int32) (Branch, error)
	GetBranchesByBrandId(ctx context.Context, brandID int32) ([]GetBranchesByBrandIdRow, error)
	SetBranchBrandGoal(ctx context.Context, arg SetBranchBrandGoalParams) error
	GetAllBranches(ctx context.Context) ([]Branch, error)
	SetBrandSaleTypeGoal(ctx context.Context, arg SetBrandSaleTypeGoalParams) error
	GetBrandOverallGoalByGivenDateRange(ctx context.Context, arg GetBrandOverallGoalByGivenDateRangeParams) (int64, error)
	AddBranch(ctx context.Context, arg AddBranchParams) (int32, error)
	AddBranchBrand(ctx context.Context, arg AddBranchBrandParams) error
	AddBranchWithBrandsTX(ctx context.Context, branch entities.BranchWithBrands) error
	UpdateBranch(ctx context.Context, arg UpdateBranchParams) error
	UpdateBranchWithBrandsTX(ctx context.Context, branch entities.BranchWithBrands) error
	DeleteBranchBrands(ctx context.Context, branchID int32) error
}

func (db *DBStore) AddBranchWithBrandsTX(ctx context.Context, branch entities.BranchWithBrands) error {
	return db.execTx(ctx, func(queries *Queries) error {
		id, err := queries.AddBranch(ctx, AddBranchParams{
			Title:       branch.Title,
			Description: branch.Description,
		})
		if err != nil {
			return err
		}
		for _, brandID := range branch.BrandIDs {
			err = queries.AddBranchBrand(ctx, AddBranchBrandParams{
				BranchID: id,
				BrandID:  brandID,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *DBStore) UpdateBranchWithBrandsTX(ctx context.Context, branch entities.BranchWithBrands) error {
	return db.execTx(ctx, func(queries *Queries) error {
		err := queries.UpdateBranch(ctx, UpdateBranchParams{
			ID:          branch.BranchId,
			Title:       branch.Title,
			Description: branch.Description,
		})
		if err != nil {
			return err
		}

		err = queries.DeleteBranchBrands(ctx, branch.BranchId)
		if err != nil {
			return err
		}

		for _, brandID := range branch.BrandIDs {
			err = queries.AddBranchBrand(ctx, AddBranchBrandParams{
				BranchID: branch.BranchId,
				BrandID:  brandID,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}
