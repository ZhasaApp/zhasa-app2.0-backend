package entities

type SalesManagerMonthStatisticRequestBody struct {
	UserId int32 `json:"user_id"`
	Month  int32 `json:"month"`
	Year   int32 `json:"year"`
}

type SalesManagerYearStatisticRequestBody struct {
	UserId int32 `json:"user_id"`
	Year   int32 `json:"year"`
}

type SaleTypeResponse struct {
	Title string `json:"title"`
	Color string `json:"color"`
}
type SalesManagerYearStatisticResponseBody struct {
	SaleType SaleTypeResponse `json:"sale_type"`
	Month    int32            `json:"month"`
	Amount   int64            `json:"amount"`
}
