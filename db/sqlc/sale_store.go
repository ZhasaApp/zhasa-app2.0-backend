package generated

import "context"

type SaleStore interface {
	AddBrandSaleTx(ctx context.Context, params AddSaleOrReplaceParams, brandId int32) (Sale, error)
	GetSaleSumByUserIdBrandIdPeriodSaleTypeId(ctx context.Context, arg GetSaleSumByUserIdBrandIdPeriodSaleTypeIdParams) (int64, error)
}
