package brand

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type GetAllBrandsFunc func() ([]Brand, error)

func NewGetAllBrandsFunc(ctx context.Context, store generated.BrandStore) GetAllBrandsFunc {
	return func() ([]Brand, error) {
		rows, err := store.GetBrands(ctx, generated.GetBrandsParams{
			Limit:  10,
			Offset: 0,
		})
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
