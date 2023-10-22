package entities

import (
	"zhasa2.0/date"
	sale "zhasa2.0/sale/entities"
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

func SaleItemsFromSales(rows []sale.Sale) []SaleItemResponse {
	result := make([]SaleItemResponse, 0)
	for _, row := range rows {
		result = append(result, SaleItemResponse{
			Id:     row.Id,
			Title:  row.SaleDescription,
			Date:   date.ConvertTimeToStringISO(row.SaleDate),
			Amount: row.SalesAmount,
			Type: SaleTypeResponse{
				Id:        row.SaleType.Id,
				Title:     row.SaleType.Title,
				Color:     row.SaleType.Color,
				ValueType: row.SaleType.ValueType,
			},
		})
	}
	return result
}
