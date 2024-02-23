package brand

import (
	"context"
	generated "zhasa2.0/db/sqlc"
)

type UpdateBrandFunc func(brand Brand) error

func NewUpdateBrandFunc(ctx context.Context, store generated.BrandStore) UpdateBrandFunc {
	return func(brand Brand) error {
		return store.UpdateBrand(ctx, generated.UpdateBrandParams{
			ID:    brand.Id,
			Title: brand.Title,
		})
	}
}
