package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/statistic"
)

type GetBranchGoalRequest struct {
	BranchId   int32 `json:"branch_id" form:"branch_id"`
	BrandId    int32 `json:"brand_id" form:"brand_id"`
	Month      int32 `json:"month" form:"month"`
	Year       int32 `json:"year" form:"year"`
	SaleTypeID int32 `json:"sale_type_id" form:"sale_type_id"`
}

type GetBranchGoalResponse struct {
	Value *int64 `json:"value"`
}

func (server *Server) GetBranchGoal(ctx *gin.Context) {
	var request GetBranchGoalRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	period := statistic.MonthPeriod{
		MonthNumber: request.Month,
		Year:        request.Year,
	}

	goal, err := server.getBranchBrandGoalFunc(request.BranchId, request.BrandId, request.SaleTypeID, period)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if goal == 0 {
		ctx.JSON(http.StatusOK, GetBranchGoalResponse{
			Value: nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, GetBranchGoalResponse{
		Value: &(goal),
	})
}
