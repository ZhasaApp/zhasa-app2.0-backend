package hand_made

import (
	"context"
	"database/sql"
)

type CustomQuerier interface {
	GetSalesManagerYearStatistic(ctx context.Context, arg GetSalesManagerYearStatisticParams) (*GetSalesManagerYearStatisticRow, error)
	GetBranchYearStatistic(ctx context.Context, arg GetBranchYearStatisticParams) (*GetBranchYearStatisticRow, error)
	GetBranchSumByType(ctx context.Context, arg GetBranchSumByTypeParams) (GetBranchSumByTypeRow, error)
	GetPostsAndPostAuthors(ctx context.Context, arg GetPostsAndPostAuthorsParams) ([]GetPostsAndPostAuthorsRow, error)
	GetSaleSumByUserIdBrandIdPeriodSaleTypeId(ctx context.Context, arg GetSaleSumByUserIdBrandIdPeriodSaleTypeIdParams) (int64, error)
}

func NewCustomQuerier(db *sql.DB) CustomQuerier {
	return DBCustomQuerier{
		db: db,
	}
}

type DBCustomQuerier struct {
	db *sql.DB
}
