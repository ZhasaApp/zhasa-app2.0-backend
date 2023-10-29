package brand

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type GetBranchBrandsFunc func(branchId int32) ([]Brand, error)

func NewGetBranchBrandsFunc(ctx context.Context, store generated.BrandStore) GetBranchBrandsFunc {
	return func(branchId int32) ([]Brand, error) {
		rows, err := store.GetBranchBrands(ctx, branchId)
		if err != nil {
			return nil, err
		}
		brands := make([]Brand, 0)
		for _, row := range rows {
			brands = append(brands, Brand{
				Id:    row.ID,
				Title: row.Title,
			})
		}

		return brands, nil
	}
}
