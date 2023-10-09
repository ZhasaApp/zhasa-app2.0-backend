package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/api/entities"
)

type GetBranchBrandsRequest struct {
	Id int32 `json:"id"`
}

func (server *Server) GetBranchBrands(ctx *gin.Context) {
	var getBranchBrandsRequest GetBranchBrandsRequest
	if err := ctx.ShouldBindJSON(&getBranchBrandsRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	brands, err := server.getBranchBrands(getBranchBrandsRequest.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := entities.BrandsResponse{
		Result: entities.BrandItemsFromBrands(brands),
	}

	ctx.JSON(http.StatusOK, response)
}
