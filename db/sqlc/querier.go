// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package generated

import (
	"context"
	"time"
)

type Querier interface {
	// add sale into sales by given sale_type_id, amount, date, sales_manager_id and on conflict replace
	AddSaleOrReplace(ctx context.Context, arg AddSaleOrReplaceParams) (Sale, error)
	CreateBranch(ctx context.Context, arg CreateBranchParams) error
	CreateBranchDirector(ctx context.Context, arg CreateBranchDirectorParams) (int32, error)
	CreateSaleType(ctx context.Context, arg CreateSaleTypeParams) (int32, error)
	CreateSalesManager(ctx context.Context, arg CreateSalesManagerParams) error
	CreateSalesManagerGoalByType(ctx context.Context, arg CreateSalesManagerGoalByTypeParams) error
	CreateUser(ctx context.Context, arg CreateUserParams) error
	CreateUserCode(ctx context.Context, arg CreateUserCodeParams) (int32, error)
	DeleteSaleById(ctx context.Context, id int32) (Sale, error)
	EditSaleById(ctx context.Context, arg EditSaleByIdParams) (Sale, error)
	GetAuthCodeById(ctx context.Context, id int32) (UsersCode, error)
	GetBranchById(ctx context.Context, id int32) (Branch, error)
	GetBranchDirectorByUserId(ctx context.Context, userID int32) (BranchDirectorsView, error)
	GetBranchGoalByGivenDateRange(ctx context.Context, arg GetBranchGoalByGivenDateRangeParams) (int64, error)
	GetManagerSales(ctx context.Context, arg GetManagerSalesParams) ([]GetManagerSalesRow, error)
	GetManagerSalesByPeriod(ctx context.Context, arg GetManagerSalesByPeriodParams) ([]GetManagerSalesByPeriodRow, error)
	GetOrderedBranchesByGivenPeriod(ctx context.Context, arg GetOrderedBranchesByGivenPeriodParams) ([]GetOrderedBranchesByGivenPeriodRow, error)
	GetOrderedSalesManagers(ctx context.Context, arg GetOrderedSalesManagersParams) ([]GetOrderedSalesManagersRow, error)
	GetOrderedSalesManagersOfBranch(ctx context.Context, arg GetOrderedSalesManagersOfBranchParams) ([]GetOrderedSalesManagersOfBranchRow, error)
	GetSMGoal(ctx context.Context, arg GetSMGoalParams) (int64, error)
	GetSMRatio(ctx context.Context, arg GetSMRatioParams) (float64, error)
	GetSaleTypeById(ctx context.Context, id int32) (SaleType, error)
	GetSalesByDate(ctx context.Context, saleDate time.Time) ([]Sale, error)
	GetSalesCount(ctx context.Context, salesManagerID int32) (int64, error)
	GetSalesManagerByUserId(ctx context.Context, userID int32) (SalesManagersView, error)
	GetSalesManagerGoalByGivenDateRangeAndSaleType(ctx context.Context, arg GetSalesManagerGoalByGivenDateRangeAndSaleTypeParams) (int64, error)
	// get the sales sums for a specific sales manager and each sale type within the given period.
	GetSalesManagerSumsByType(ctx context.Context, arg GetSalesManagerSumsByTypeParams) (GetSalesManagerSumsByTypeRow, error)
	GetSalesTypes(ctx context.Context) ([]SaleType, error)
	GetUserById(ctx context.Context, id int32) (User, error)
	GetUserByPhone(ctx context.Context, phone string) (User, error)
	SetSMRatio(ctx context.Context, arg SetSMRatioParams) error
	SetSmGoalBySaleType(ctx context.Context, arg SetSmGoalBySaleTypeParams) error
}

var _ Querier = (*Queries)(nil)
