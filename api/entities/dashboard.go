package entities

type CreateSalesManagerBody struct {
	CreateUserBody
	BranchId int32 `json:"branch_id"`
}

type SaveSaleBody struct {
	SaleAmount  int64  `json:"amount"`
	SaleDate    string `json:"date"`
	SaleTypeId  int32  `json:"type_id"`
	Description string `json:"description"`
}

type OverallSalesStatistic struct {
	Goal         int64        `json:"goal"`
	Achieved     int64        `json:"achieved"`
	Percent      float64      `json:"percent"`
	GrowthPerDay GrowthPerDay `json:"growth_per_day"`
}

type GrowthPerDay struct {
	Amount  int64   `json:"amount"`
	Percent float64 `json:"percent"`
}

type SalesManagerDashboardResponse struct {
	OverallSalesStatistics OverallSalesStatistic        `json:"overall_sales_statistics"`
	SalesStatisticsByTypes []SalesStatisticsByTypesItem `json:"sales_statistics_by_types"`
	LastSales              []SaleItemResponse           `json:"last_sales"`
}

type SalesManagerBranchItem struct {
	Id          int32   `json:"id"`
	Avatar      string  `json:"avatar"`
	FullName    string  `json:"full_name"`
	Ratio       float64 `json:"ratio"`
	BranchTitle string  `json:"branch_title"`
	BranchId    int32   `json:"branch_id"`
}

type BranchDashboardResponse struct {
	OverallSaleStatistics OverallSalesStatistic        `json:"overall_sale_statistics"`
	SaleStatisticsByTypes []SalesStatisticsByTypesItem `json:"sale_statistics_by_types"`
	BestSalesManagers     []SalesManagerBranchItem     `json:"best_sales_managers"`
}

type SaleItemResponse struct {
	Id     int32            `json:"id"`
	Title  string           `json:"title"`
	Date   string           `json:"date"`
	Amount int64            `json:"amount"`
	Type   SaleTypeResponse `json:"type"`
}

type SalesStatisticsByTypesItem struct {
	Color  string `json:"color"`
	Title  string `json:"title"`
	Amount int64  `json:"amount"`
}
