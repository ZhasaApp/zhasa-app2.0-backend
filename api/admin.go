package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/sale/entities"
)

type createSaleTypeBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type createSaleTypeResponse struct {
	Id int32 `json:"id"`
}

func (server *Server) createSaleType(ctx *gin.Context) {
	var createSaleTypeBody createSaleTypeBody
	if err := ctx.ShouldBindJSON(&createSaleTypeBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	body := entities.CreateSaleTypeBody{
		Title:       createSaleTypeBody.Title,
		Description: createSaleTypeBody.Description,
	}

	id, err := server.saleTypeService.CreateSaleType(body)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, createSaleTypeResponse{
		Id: int32(id),
	})
}
