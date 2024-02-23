package brand

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type CreateBrandFunc func(brand Brand) error

func NewCreateBrandFunc(ctx context.Context, store generated.BrandStore) CreateBrandFunc {
	return func(brand Brand) error {
		return store.AddBrand(ctx, generated.AddBrandParams{
			Title: brand.Title,
		})
	}
}
