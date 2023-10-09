package brand

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type GetUserBrandsFunc func(userId int32) ([]Brand, error)

func NewGetUserBrandsFunc(ctx context.Context, store generated.BrandStore) GetUserBrandsFunc {
	return func(userId int32) ([]Brand, error) {
		branchBrands, err := store.GetUserBrands(ctx, userId)
		if err != nil {
			return nil, err
		}

		brands := make([]Brand, 0)
		for _, branchBrand := range branchBrands {
			brands = append(brands, Brand{
				Id:    branchBrand.ID,
				Title: branchBrand.Title,
			})
		}

		return brands, nil
	}
}
