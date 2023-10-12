package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	entities2 "zhasa2.0/api/entities"
	"zhasa2.0/base"
	"zhasa2.0/rating"
	"zhasa2.0/statistic"
)

type GetBranchDashboardRequest struct {
	BranchId int32 `json:"branch_id" form:"branch_id"`
	BrandId  int32 `json:"brand_id" form:"brand_id"`
	Month    int32 `json:"month" form:"month"`
	Year     int32 `json:"year" form:"year"`
}

func (server *Server) BranchDashboard(ctx *gin.Context) {
	var request GetBranchDashboardRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	branchInfo, err := server.getBranchByIdFunc(request.BranchId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.New("branch not found"))
		return
	}

	branchBrand, err := server.getBranchBrandFunc(request.BranchId, request.BrandId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.New("branch brand not found"))
		return
	}

	saleTypes, err := server.saleTypeRepo.GetSaleTypes()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.New("sale types not found"))
		return
	}

	monthPeriod := statistic.MonthPeriod{
		MonthNumber: request.Month,
		Year:        request.Year,
	}

	var goalAchievementPercent float32
	ratioRows := make([]rating.RatioRow, 0)
	saleStatisticByTypes := make([]entities2.SalesStatisticsByTypesItem, 0)

	for _, saleType := range *saleTypes {
		achieved, _ := server.getBranchBrandSaleSumFunc(request.BranchId, request.BrandId, saleType.Id, monthPeriod)
		goal, _ := server.getBranchBrandGoalFunc(branchBrand, saleType.Id, monthPeriod)

		ratioRows = append(ratioRows, rating.RatioRow{
			Achieved: achieved,
			Goal:     goal,
			Gravity:  saleType.Gravity,
		})

		saleStatisticByTypes = append(saleStatisticByTypes, entities2.SalesStatisticsByTypesItem{
			Color:    saleType.Color,
			Title:    saleType.Title,
			Achieved: achieved,
			Goal:     goal,
		})

	}

	goalAchievementPercent = rating.CalculateRatio(ratioRows)

	bestSalesManagers, err := server.getBranchUsersOrderedByRatioForGivenBrandFunc(request.BrandId, request.BranchId, monthPeriod, base.Pagination{
		PageSize: 3,
		Page:     0,
	})

	ctx.JSON(http.StatusOK, entities2.BranchDashboardResponse{
		SalesStatisticsByTypes: saleStatisticByTypes,
		BestSalesManagers:      entities2.SalesManagerBranchItemsFromRatedUsers(bestSalesManagers),
		Branch: entities2.BranchModelResponse{
			Title:       branchInfo.Title,
			Description: branchInfo.Description,
		},
		GoalAchievementPercent: goalAchievementPercent,
		Rating:                 0,
		Profile:                entities2.SimpleProfile{},
	})

}
