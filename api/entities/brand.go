package entities

import generated "zhasa2.0/db/sqlc"

type Brand struct {
	Id    int32  `json:"id"`
	Title string `json:"title"`
}

type BrandsResponse struct {
	Result []Brand `json:"result"`
}

func BrandsFromRows(rows []generated.GetBranchBrandsRow) []Brand {
	var brands []Brand
	for _, row := range rows {
		brands = append(brands, BrandFromRow(row))
	}
	return brands
}

func BrandFromRow(row generated.GetBranchBrandsRow) Brand {
	return Brand{
		Id:    row.ID,
		Title: row.Title,
	}
}
