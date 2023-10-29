package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/statistic"
)

type SetBranchGoalRequest struct {
	BranchId   int32 `json:"branch_id"`
	BrandId    int32 `json:"brand_id"`
	Value      int64 `json:"value"`
	Month      int32 `json:"month"`
	Year       int32 `json:"year"`
	SaleTypeID int32 `json:"sale_type_id"`
}

func (server *Server) SetBranchGoal(ctx *gin.Context) {
	var request SetBranchGoalRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	period := statistic.MonthPeriod{
		MonthNumber: request.Month,
		Year:        request.Year,
	}

	err := server.setBranchBrandSaleTypeGoal(request.BranchId, request.BrandId, request.SaleTypeID, request.Value, period)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}
