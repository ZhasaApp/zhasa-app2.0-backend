package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/api/entities"
	"zhasa2.0/base"
)

type GetSalesRequest struct {
	UserId  int32 `json:"user_id" form:"user_id"`
	BrandId int32 `json:"brand_id" form:"brand_id"`
	Month   int32 `json:"month" form:"month"`
	Year    int32 `json:"year" form:"year"`
	Page    int32 `json:"page" form:"page"`
	Limit   int32 `json:"limit" form:"limit"`
}

func (server *Server) GetSales(ctx *gin.Context) {
	var request GetSalesRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fmt.Println(request)

	sales, err := server.salesByBrandUserFunc(request.UserId, request.BrandId, base.Pagination{
		PageSize: request.Limit,
		Page:     request.Page,
	})

	nextSales, err := server.salesByBrandUserFunc(request.UserId, request.BrandId, base.Pagination{
		PageSize: request.Limit,
		Page:     request.Page + 1,
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
