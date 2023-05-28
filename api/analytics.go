package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	. "zhasa2.0/api/entities"
	. "zhasa2.0/base"
	. "zhasa2.0/statistic/entities"
)

func (server *Server) getRankedSalesManager(ctx *gin.Context) {
	var requestBody MonthPaginationRequest
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	period := MonthPeriod{
		MonthNumber: requestBody.Month,
		Year:        requestBody.Year,
	}

	pagination := Pagination{
		PageSize: requestBody.PageSize,
		Page:     requestBody.Page,
	}

	managers, err := server.analyticsService.ProvideRankedManagers(pagination, period)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	bestSalesManagers := make([]SalesManagerBranchItem, 0)

	for _, sm := range *managers {
		bestSalesManagers = append(bestSalesManagers, SalesManagerBranchItem{
			Id:          int32(sm.UserId),
			Avatar:      sm.AvatarUrl,
			FullName:    sm.FirstName + " " + sm.LastName,
			Ratio:       float64(sm.Ratio),
			BranchTitle: string(sm.Branch.Title),
			BranchId:    int32(sm.Branch.BranchId),
		})
	}

	ctx.JSON(http.StatusOK, bestSalesManagers)
}
