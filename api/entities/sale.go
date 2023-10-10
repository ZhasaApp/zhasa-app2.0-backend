package entities

import (
	"zhasa2.0/date"
	generated "zhasa2.0/db/sqlc"
)

type SaleItemResponse struct {
	Id     int32            `json:"id"`
	Title  string           `json:"title"`
	Date   string           `json:"date"`
	Amount int64            `json:"value"`
	Type   SaleTypeResponse `json:"type"`
}

type SalesResponse struct {
	Result  []SaleItemResponse `json:"result"`
	Count   int32              `json:"count"`
	HasNext bool               `json:"has_next"`
}

func SaleItemsFromSales(rows []generated.GetSalesByBrandIdAndUserIdRow) []SaleItemResponse {
	result := make([]SaleItemResponse, 0)
	for _, row := range rows {
		result = append(result, SaleItemResponse{
			Id:     row.ID,
			Title:  row.Description,
			Date:   date.ConvertTimeToStringISO(row.SaleDate),
			Amount: row.Amount,
			Type: SaleTypeResponse{
				Id:        row.SaleTypeID,
				Title:     row.SaleTypeTitle,
				Color:     row.Color,
				ValueType: string(row.ValueType),
			},
		})
	}
	return result
}
