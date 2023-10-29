package entities

import (
	"zhasa2.0/brand"
)

type BrandItem struct {
	Id    int32  `json:"id"`
	Title string `json:"title"`
}

type BrandsResponse struct {
	Result []BrandItem `json:"result"`
}

func BrandItemsFromBrands(rows []brand.Brand) []BrandItem {
	var brands []BrandItem
	for _, row := range rows {
		brands = append(brands, BrandFromRow(row))
	}
	return brands
}

func BrandFromRow(row brand.Brand) BrandItem {
	return BrandItem{
		Id:    row.Id,
		Title: row.Title,
	}
}
