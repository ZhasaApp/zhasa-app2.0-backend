package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type SetBranchGoalRequest struct {
	BranchId   int32 `json:"branch_id"`
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
	//
	//period := entities.MonthPeriod{
	//	MonthNumber: request.Month,
	//	Year:        request.Year,
	//}
	//
	////err := server.directorRepo.SetBranchGoal(period, request.BranchId, request.SaleTypeID, request.Value)
	//
	//if err != nil {
	//	ctx.JSON(http.StatusBadRequest, err)
	//	return
	//}

	ctx.Status(http.StatusNoContent)
}
