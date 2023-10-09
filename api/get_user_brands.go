package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/api/entities"
)

type GetUserBrandsRequest struct {
	Id int32 `json:"id" form:"id" binding:"required"`
}

func (server *Server) GetUserBrands(ctx *gin.Context) {
	var getUserBrandsRequest GetUserBrandsRequest
	if err := ctx.ShouldBindQuery(&getUserBrandsRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	brands, err := server.getUserBrands(getUserBrandsRequest.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := entities.BrandsResponse{
		Result: entities.BrandItemsFromBrands(brands),
	}

	ctx.JSON(http.StatusOK, response)
}
