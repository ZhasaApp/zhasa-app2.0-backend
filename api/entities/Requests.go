package entities

type SalesManagerMonthStatisticRequestBody struct {
	UserId int32 `json:"user_id"`
	Month  int32 `json:"month"`
	Year   int32 `json:"year"`
}
