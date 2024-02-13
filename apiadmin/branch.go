package apiadmin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/branch/entities"
)

type CreateBranchWithBrandsRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	BrandIDs    []int32 `json:"brand_ids"`
}

func (s *Server) CreateBranchWithBrands(ctx *gin.Context) {
	var request CreateBranchWithBrandsRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.createBranchWithBrandsFunc(entities.BranchWithBrands{
		Branch: entities.Branch{
			Title:       request.Title,
			Description: request.Description,
		},
		BrandIDs: request.BrandIDs,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.Status(http.StatusCreated)
}
