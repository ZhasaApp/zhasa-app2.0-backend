package entities

type CreateSaleTypeBody struct {
	Title       string
	Description string
}

type SaleType struct {
	Id          int32
	Title       string
	Description string
	Color       string
	Gravity     int32
	ValueType   string
}

type SumsByTypeRow struct {
	SaleTypeID    int32  `json:"sale_type_id"`
	SaleTypeTitle string `json:"sale_type_title"`
	TotalSales    int64  `json:"total_sales"`
}
