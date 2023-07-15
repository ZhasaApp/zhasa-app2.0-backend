package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/statistic/entities"
)

type SetGoalRequest struct {
	UserID     int32 `json:"user_id"`
	Value      int64 `json:"value"`
	Month      int32 `json:"month"`
	Year       int32 `json:"year"`
	SaleTypeID int32 `json:"sale_type_id"`
}

func (server *Server) SetSmGoal(ctx *gin.Context) {
	var request SetGoalRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
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

	err = server.directorService.SetSmGoal(period, int32(sm.Id), request.SaleTypeID, request.Value)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
