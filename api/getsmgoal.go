package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/statistic"
)

type GetGoalRequest struct {
	UserID     int32 `json:"user_id" form:"user_id"`
	BrandId    int32 `json:"brand_id" form:"brand_id"`
	Month      int32 `json:"month" form:"month"`
	Year       int32 `json:"year" form:"year"`
	SaleTypeID int32 `json:"sale_type_id" form:"sale_type_id"`
}

type GetGoalResponse struct {
	Value *int64 `json:"value"`
}

func (server *Server) GetSmGoal(ctx *gin.Context) {
	var request GetGoalRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	period := statistic.MonthPeriod{
		MonthNumber: request.Month,
		Year:        request.Year,
	}

	goal := server.userBrandGoal(request.UserID, request.BrandId, request.SaleTypeID, period)

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
