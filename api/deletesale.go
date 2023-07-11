package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	. "zhasa2.0/sale/entities"
	. "zhasa2.0/statistic/entities"
)

type DeleteSaleQuery struct {
	SaleId int32 `json:"id" form:"id"`
}

func (server Server) DeleteSale(ctx *gin.Context) {
	var requestBody DeleteSaleQuery
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	deletedSale, err := server.salesManagerService.DeleteSale(SaleId(requestBody.SaleId))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	server.salesManagerService.UpdateRatio(deletedSale.SaleManagerId, MonthPeriod{
		MonthNumber: int32(deletedSale.SaleDate.Month()),
		Year:        int32(deletedSale.SaleDate.Year()),
	})

	ctx.Status(http.StatusNoContent)
}
