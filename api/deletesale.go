package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	. "zhasa2.0/sale/entities"
)

type DeleteSaleQuery struct {
	SaleId int32 `json:"id" form:"id"`
}

func (server Server) DeleteSale(ctx *gin.Context) {
	var requestBody DeleteSaleQuery
	if err := ctx.ShouldBindQuery(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.salesManagerService.DeleteSale(SaleId(requestBody.SaleId))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}
