package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetGoodBySaleIdRequest struct {
	SaleId int32 `json:"sale_id" form:"sale_id" binding:"required"`
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
		ctx.JSON(http.StatusNotFound, errorResponse(errors.New("good not found by sale id")))
		return
	}

	response := GoodItem{
		Id:          good.Id,
		Name:        good.Name,
		Description: good.Description,
	}

	ctx.JSON(http.StatusOK, response)
}
