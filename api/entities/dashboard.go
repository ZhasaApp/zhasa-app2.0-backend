package entities

type CreateSalesManagerBody struct {
	CreateUserBody
	BranchId int32 `json:"branch_id"`
}

type SaveSaleBody struct {
	SaleAmount  int64  `json:"sale_amount"`
	SaleDate    string `json:"sale_date"`
	SaleTypeId  int32  `json:"sale_type_id"`
	Description string `json:"description"`
}

type OverallSaleStatistic struct {
	Goal         int64        `json:"goal"`
	Achieved     int64        `json:"achieved"`
	Percent      float64      `json:"percent"`
	GrowthPerDay GrowthPerDay `json:"growth_per_day"`
}

type GrowthPerDay struct {
	Amount  int64   `json:"amount"`
	Percent float64 `json:"percent"`
}

type DashboardResponse struct {
	OverallSaleStatistics OverallSaleStatistic        `json:"overall_sale_statistics"`
	SaleStatisticsByTypes []SaleStatisticsByTypesItem `json:"sale_statistics_by_types"`
	LastSales             []SaleItemResponse          `json:"last_sales"`
}

type SaleItemResponse struct {
	Id          int32  `json:"id"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Amount      int64  `json:"amount"`
}

type SaleStatisticsByTypesItem struct {
	Color  string `json:"color"`
	Title  string `json:"title"`
	Amount int64  `json:"amount"`
}
