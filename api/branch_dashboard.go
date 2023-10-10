package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/api/entities"
)

type GetBranchDashboardRequest struct {
	BranchId int32 `json:"branch_id" form:"branch_id"`
	BrandId  int32 `json:"brand_id" form:"brand_id"`
	Month    int32 `json:"month" form:"month"`
	Year     int32 `json:"year" form:"year"`
}

func (server *Server) BranchDashboard(ctx *gin.Context) {
	var request GetSMDashboardRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entities.BranchDashboardResponse{
		SalesStatisticsByTypes: nil,
		BestSalesManagers:      nil,
		Branch:                 entities.BranchModelResponse{},
		GoalAchievementPercent: 0,
		Rating:                 0,
		Profile:                entities.SimpleProfile{},
	})
}
