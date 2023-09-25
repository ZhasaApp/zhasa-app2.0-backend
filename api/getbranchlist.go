package api

import (
	"github.com/gin-gonic/gin"
)

type BranchRatingItem struct {
	ID                     int32   `json:"id"`
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
	//var monthPagination MonthPaginationRequest
	//if err := ctx.ShouldBindQuery(&monthPagination); err != nil {
	//	ctx.JSON(http.StatusBadRequest, errorResponse(err))
	//	return
	//}
	//
	//branchesResponseList := make([]BranchRatingItem, 0)
	//
	//branches, err := server.branchService.GetBranches(MonthPeriod{
	//	MonthNumber: monthPagination.Month,
	//	Year:        monthPagination.Year,
	//})
	//
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	//	return
	//}
	//
	//for _, branch := range branches {
	//	branchesResponseList = append(branchesResponseList, BranchRatingItem{
	//		ID:                     int32(branch.BranchId),
	//		Title:                  string(branch.Title),
	//		Description:            string(branch.Description),
	//		GoalAchievementPercent: float32(branch.GoalAchievement),
	//	})
	//}
	//sort.Slice(branchesResponseList, func(i, j int) bool {
	//	return branchesResponseList[i].GoalAchievementPercent > branchesResponseList[j].GoalAchievementPercent
	//})
	//ctx.JSON(http.StatusOK, BranchesResponse{
	//	Result:  branchesResponseList,
	//	Count:   int32(len(branchesResponseList)),
	//	HasNext: false,
	//})
}
