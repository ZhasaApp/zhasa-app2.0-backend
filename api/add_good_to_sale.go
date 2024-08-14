package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/sale/repository"
)

type AddGoodToSaleRequest struct {
	SaleId int32 `json:"sale_id"`
	GoodId int32 `json:"good_id"`
}

func (server *Server) AddGoodToSale(ctx *gin.Context) {
	var requestBody AddGoodToSaleRequest
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.saleAddWithGoodFunc(repository.AddGoodToSaleBody{
		SaleId: requestBody.SaleId,
		GoodId: requestBody.GoodId,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}
