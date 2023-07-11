package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	. "zhasa2.0/api/entities"
	. "zhasa2.0/statistic/entities"
)

type BranchRatingItem struct {
	ID                     int     `json:"id"`
	Title                  string  `json:"title"`
	Description            string  `json:"description"`
	GoalAchievementPercent float32 `json:"goal_achievement_percent"`
}

type BranchesResponse struct {
	Result  []BranchRatingItem `json:"result"`
	Count   int32              `json:"count"`
	HasNext bool               `json:"has_next"`
}

func (server *Server) GetBranchList(ctx *gin.Context) {
	var monthPagination MonthPaginationRequest
	if err := ctx.ShouldBindQuery(&monthPagination); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, BranchesResponse{})
}
