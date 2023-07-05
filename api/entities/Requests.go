package entities

type SalesManagerMonthStatisticRequestBody struct {
	UserId int32 `form:"user_id" json:"user_id" binding:"required"`
	Month  int32 `form:"month" json:"month" binding:"required"`
	Year   int32 `form:"year" json:"year" binding:"required"`
}

type BranchMonthStatisticRequestBody struct {
	BranchId int32 `json:"branch_id"`
	Month    int32 `json:"month"`
	Year     int32 `json:"year"`
}

type SalesManagerYearStatisticRequestBody struct {
	UserId int32 `form:"user_id" json:"user_id"`
	Year   int32 `form:"year" json:"year"`
}

type BranchYearStatisticRequestBody struct {
	BranchId int32 `json:"branch_id"`
	Year     int32 `json:"year"`
}

type SaleTypeResponse struct {
	Id    int32  `json:"id"`
	Title string `json:"title"`
	Color string `json:"color"`
}

type SaleTypesResponse struct {
	Result []SaleTypeResponse `json:"result"`
}

type YearStatisticResultResponse struct {
	Result []YearStatisticResponse `json:"result"`
}

type YearStatisticResponse struct {
	SaleType SaleTypeResponse `json:"sale_type"`
	Month    int32            `json:"month"`
	Amount   int64            `json:"value"`
	Goal     int64            `json:"goal"`
}

type MonthPaginationRequest struct {
	Month    int32 `json:"month" form:"month"`
	Year     int32 `json:"year" form:"year"`
	Page     int32 `json:"page" form:"page"`
	PageSize int32 `json:"limit" form:"limit"`
	UserId   int32 `json:"user_id" form:"user_id"`
	BranchId int32 `json:"branchId" form:"branch_id"`
}
