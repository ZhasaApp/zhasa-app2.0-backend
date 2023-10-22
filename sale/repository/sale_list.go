package repository

import (
	"context"
	"database/sql"
	"zhasa2.0/base"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/sale/entities"
)

type SalesByBrandUserFunc func(userId, brandId int32, pagination base.Pagination) ([]entities.Sale, error)

func NewSalesByBrandUserFunc(ctx context.Context, store generated.SaleStore) SalesByBrandUserFunc {
	return func(userId, brandId int32, pagination base.Pagination) ([]entities.Sale, error) {
		rows, err := store.GetSalesByBrandIdAndUserId(ctx, generated.GetSalesByBrandIdAndUserIdParams{
			BrandID: brandId,
			UserID:  userId,
			Limit:   pagination.PageSize,
			Offset:  pagination.GetOffset(),
		})

		sales := make([]entities.Sale, 0)
		if err == sql.ErrNoRows {
			return sales, nil
		}

		if err != nil {
			return nil, err
		}

		for _, row := range rows {
			sales = append(sales, entities.Sale{
				Id: row.ID,
				SaleType: entities.SaleType{
					Id:          row.SaleTypeID,
					Title:       row.SaleTypeTitle,
					Description: "",
					Color:       row.Color,
					Gravity:     row.Gravity,
					ValueType:   string(row.ValueType),
				},
				SalesAmount:     row.Amount,
				SaleDate:        row.SaleDate,
				SaleDescription: row.Description,
			})
		}
		return sales, nil
	}
}
