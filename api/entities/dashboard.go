package entities

type CreateSalesManagerBody struct {
	CreateUserBody
	BranchId int32 `json:"branch_id"`
}

type SaveSaleBody struct {
	SaleAmount  int64  `json:"value"`
	SaleDate    string `json:"date"`
	SaleTypeId  int32  `json:"type_id"`
	Description string `json:"title"`
}

type GrowthPerDay struct {
	Amount  int64   `json:"amount"`
	Percent float64 `json:"percent"`
}

type SalesManagerDashboardResponse struct {
	Profile                SalesManagerDashboardProfile `json:"profile"`
	SalesStatisticsByTypes []SalesStatisticsByTypesItem `json:"sales_statistics_by_types"`
	GoalAchievementPercent float32                      `json:"goal_achievement_percent"`
	LastSales              []SaleItemResponse           `json:"last_sales"`
	Rating                 int32                        `json:"rating"`
}

type SalesManagerBranchItem struct {
	Id          int32   `json:"id"`
	Avatar      *string `json:"avatar"`
	FullName    string  `json:"full_name"`
	Ratio       float64 `json:"goal_achievement_percent"`
	BranchTitle string  `json:"branch"`
	BranchId    int32   `json:"branch_id"`
}

type SalesManagersListResponse struct {
	Result  []SalesManagerBranchItem `json:"result"`
	Count   int32                    `json:"count"`
	HasNext bool
}

type BranchModelResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type BranchDashboardResponse struct {
	SalesStatisticsByTypes []SalesStatisticsByTypesItem `json:"sales_statistics_by_types"`
	BestSalesManagers      []SalesManagerBranchItem     `json:"best_sales_managers"`
	Branch                 BranchModelResponse          `json:"branch"`
	GoalAchievementPercent float32                      `json:"goal_achievement_percent"`
	Rating                 int32                        `json:"rating"`
	Profile                SimpleProfile                `json:"profile"`
}

type SaleItemResponse struct {
	Id     int32            `json:"id"`
	Title  string           `json:"title"`
	Date   string           `json:"date"`
	Amount int64            `json:"value"`
	Type   SaleTypeResponse `json:"type"`
}

type CreateSaleResponse struct {
	Id     int32  `json:"id"`
	Title  string `json:"title"`
	Date   string `json:"date"`
	Amount int64  `json:"amount"`
	TypeId int32  `json:"type_id"`
}

type SalesResponse struct {
	Result  []SaleItemResponse `json:"result"`
	Count   int32              `json:"count"`
	HasNext bool               `json:"has_next"`
}

type SalesStatisticsByTypesItem struct {
	Color    string `json:"color"`
	Title    string `json:"title"`
	Achieved int64  `json:"achieved"`
	Goal     int64  `json:"goal"`
}

type SalesManagerDashboardProfile struct {
	Avatar   *string `json:"avatar"`
	FullName string  `json:"full_name"`
	Branch   string  `json:"branch"`
}

type SimpleProfile struct {
	Avatar   *string `json:"avatar"`
	FullName string  `json:"full_name"`
}
