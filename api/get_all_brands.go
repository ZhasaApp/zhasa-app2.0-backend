package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/api/entities"
)

func (server *Server) GetAllBrands(ctx *gin.Context) {
	brands, err := server.getAllBrands()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	response := entities.BrandsResponse{
		Result: entities.BrandItemsFromBrands(brands),
	}

	ctx.JSON(http.StatusOK, response)
}
