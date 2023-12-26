package generated

import (
	"context"
)

type SaleStore interface {
	AddBrandSaleTx(ctx context.Context, params AddSaleOrReplaceParams, brandId int32) (*Sale, error)
	DeleteSale(ctx context.Context, id int32) error
	GetSalesByBrandIdAndUserId(ctx context.Context, arg GetSalesByBrandIdAndUserIdParams) ([]GetSalesByBrandIdAndUserIdRow, error)
	GetSaleBrandBySaleId(ctx context.Context, saleID int32) (GetSaleBrandBySaleIdRow, error)
	GetSaleSumByBranchByTypeByBrand(ctx context.Context, arg GetSaleSumByBranchByTypeByBrandParams) (GetSaleSumByBranchByTypeByBrandRow, error)
	GetBranchBrandSaleSumByGivenDateRange(ctx context.Context, arg GetBranchBrandSaleSumByGivenDateRangeParams) (int64, error)
	GetSumByUserIdBrandIdPeriodSaleTypeId(ctx context.Context, arg GetSumByUserIdBrandIdPeriodSaleTypeIdParams) (int64, error)
	EditSale(ctx context.Context, arg EditSaleParams) (Sale, error)
	GetBrandSaleSumByGivenDateRange(ctx context.Context, arg GetBrandSaleSumByGivenDateRangeParams) (int64, error)
}
