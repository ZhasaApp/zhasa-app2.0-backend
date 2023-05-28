package entities

type SalesManagerMonthStatisticRequestBody struct {
	UserId int32 `json:"user_id"`
	Month  int32 `json:"month"`
	Year   int32 `json:"year"`
}

type BranchMonthStatisticRequestBody struct {
	BranchId int32 `json:"branch_id"`
	Month    int32 `json:"month"`
	Year     int32 `json:"year"`
}

type SalesManagerYearStatisticRequestBody struct {
	UserId int32 `json:"user_id"`
	Year   int32 `json:"year"`
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
type YearStatisticResponse struct {
	SaleType SaleTypeResponse `json:"sale_type"`
	Month    int32            `json:"month"`
	Amount   int64            `json:"amount"`
}

type MonthPaginationRequest struct {
	Month    int32 `json:"month"`
	Year     int32 `json:"year"`
	Page     int32 `json:"page"`
	PageSize int32 `json:"page_size"`
}
