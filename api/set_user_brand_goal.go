package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/statistic"
)

type SetUserBrandGoalRequest struct {
	UserID     int32 `json:"user_id"`
	BrandId    int32 `json:"brand_id"`
	Value      int64 `json:"value"`
	Month      int32 `json:"month"`
	Year       int32 `json:"year"`
	SaleTypeID int32 `json:"sale_type_id"`
}

func (server *Server) SetUserBrandGoal(ctx *gin.Context) {
	var request SetUserBrandGoalRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	period := statistic.MonthPeriod{
		MonthNumber: request.Month,
		Year:        request.Year,
	}

	userBrand, err := server.getUserBrandFunc(request.UserID, request.BrandId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.setUserBrandGoalRequest(userBrand, request.SaleTypeID, request.Value, period)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}
