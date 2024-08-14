package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetGoodBySaleIdRequest struct {
	SaleId int32 `json:"sale_id" form:"sale_id" binding:"required"`
}

type GoodResponse struct {
	Context *GoodItem `json:"context"`
}

func (server *Server) GetGoodBySaleId(ctx *gin.Context) {
	var getGoodBySaleIdRequest GetGoodBySaleIdRequest
	if err := ctx.ShouldBindQuery(&getGoodBySaleIdRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	good, err := server.getGoodBySaleIdFunc(getGoodBySaleIdRequest.SaleId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if good == nil {
		ctx.JSON(http.StatusOK, GoodResponse{Context: nil})
		return
	}

	response := GoodResponse{Context: &GoodItem{
		Id:          good.Id,
		Name:        good.Name,
		Description: good.Description,
	}}

	ctx.JSON(http.StatusOK, response)
}
