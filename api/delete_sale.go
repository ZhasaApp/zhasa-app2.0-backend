package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type DeleteSaleQuery struct {
	SaleId int32 `json:"id" form:"id"`
}

func (server *Server) DeleteSale(ctx *gin.Context) {
	var requestBody DeleteSaleQuery
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.saleRepo.DeleteSale(requestBody.SaleId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}
