package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/api/entities"
	generated "zhasa2.0/db/sqlc"
)

type GetSalesRequest struct {
	UserId  int32 `json:"user_id" form:"user_id" binding:"required"`
	BrandId int32 `json:"brand_id" form:"brand_id" binding:"required"`
	Month   int32 `json:"month" form:"month" binding:"required"`
	Year    int32 `json:"year" form:"year" binding:"required"`
	Page    int32 `json:"page" form:"page" binding:"required"`
	Limit   int32 `json:"limit" form:"limit" binding:"required"`
}

func (server *Server) GetSales(ctx *gin.Context) {
	var request GetSalesRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	sales, err := server.saleRepo.GetSalesByBrandIdAndUserId(generated.GetSalesByBrandIdAndUserIdParams{
		BrandID: request.BrandId,
		UserID:  request.UserId,
		Limit:   request.Limit,
		Offset:  (request.Page - 1) * request.Limit,
	})

	nextSales, err := server.saleRepo.GetSalesByBrandIdAndUserId(generated.GetSalesByBrandIdAndUserIdParams{
		BrandID: request.BrandId,
		UserID:  request.UserId,
		Limit:   request.Limit,
		Offset:  (request.Page) * request.Limit,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var hasNext bool

	hasNext = len(nextSales) == int(request.Limit)

	response := entities.SalesResponse{
		Result:  entities.SaleItemsFromSales(sales),
		Count:   int32(len(sales)),
		HasNext: hasNext,
	}

	ctx.JSON(http.StatusOK, response)
}
