package brand

import (
	"context"
	"errors"
	"fmt"
	generated "zhasa2.0/db/sqlc"
)

type GetBranchBrandFunc func(branchId int32, brandId int32) (int32, error)

func NewGetBranchBrand(ctx context.Context, branchStore generated.BranchStore) GetBranchBrandFunc {
	return func(branchId int32, brandId int32) (int32, error) {
		branchBrand, err := branchStore.GetBranchBrand(ctx, generated.GetBranchBrandParams{
			BranchID: branchId,
			BrandID:  brandId,
		})
		if err != nil {
			fmt.Println(err)
			return 0, errors.New("branch brand not found")
		}
		return branchBrand, nil
	}
}
