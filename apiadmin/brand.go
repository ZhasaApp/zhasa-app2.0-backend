package apiadmin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/brand"
)

type CreateBrandRequest struct {
	Title string `json:"title"`
}

func (s *Server) CreateBrand(ctx *gin.Context) {
	var request CreateBrandRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.createBrandFunc(brand.Brand{
		Title: request.Title,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.Status(http.StatusCreated)
}

type UpdateBrandRequest struct {
	BrandID int32  `json:"brand_id"`
	Title   string `json:"title"`
}

func (s *Server) UpdateBrand(ctx *gin.Context) {
	var request UpdateBrandRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.updateBrandFunc(brand.Brand{
		Id:    request.BrandID,
		Title: request.Title,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}
