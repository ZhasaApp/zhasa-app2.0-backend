package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	. "zhasa2.0/sale/entities"
	"zhasa2.0/statistic/entities"
)

type GetGoalRequest struct {
	UserID     int32 `json:"user_id" form:"user_id"`
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

	sm, err := server.salesManagerService.GetSalesManagerByUserId(request.UserID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("sales manager not found")))
		return
	}

	period := entities.MonthPeriod{
		MonthNumber: request.Month,
		Year:        request.Year,
	}

	goal, err := server.salesManagerService.GetSalesManagerGoalByType(period, sm.Id, SaleTypeId(request.SaleTypeID))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if goal == 0 {
		ctx.JSON(http.StatusOK, GetGoalResponse{
			Value: nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, GetGoalResponse{
		Value: (*int64)(&(goal)),
	})
}
