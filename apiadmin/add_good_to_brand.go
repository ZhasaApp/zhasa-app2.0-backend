package apiadmin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type AddGoodToBrandRequest struct {
	GoodId  int32 `json:"good_id"`
	BrandId int32 `json:"brand_id"`
}

func (s *Server) AddGoodToBrand(ctx *gin.Context) {
	var request AddGoodToBrandRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.addGoodToBrandFunc(request.GoodId, request.BrandId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}
