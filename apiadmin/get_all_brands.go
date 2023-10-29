package apiadmin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BrandItem struct {
	Id    int32  `json:"id"`
	Title string `json:"title"`
}

type BrandsResponse struct {
	Result []BrandItem `json:"result"`
}

func (s *Server) GetAllBrands(ctx *gin.Context) {
	brands, err := s.getAllBrandsFunc()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var response BrandsResponse
	for _, brand := range brands {
		response.Result = append(response.Result, BrandItem{
			Id:    brand.Id,
			Title: brand.Title,
		})
	}
	ctx.JSON(http.StatusOK, response)
}
