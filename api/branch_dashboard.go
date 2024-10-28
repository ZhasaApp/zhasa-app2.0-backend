package api

import (
	"errors"
	"fmt"
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("branch not found")))
		return
	}

	saleTypes, err := server.saleTypeRepo.GetSaleTypes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("sale types not found")))
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
		goal, _ := server.getBranchBrandGoalFunc(request.BranchId, request.BrandId, saleType.Id, monthPeriod)

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

	director, err := server.getUserByBranchBrandRoleFunc(request.BranchId, request.BrandId, 3)

	if director != nil && len(director) == 0 {
		fmt.Println("director not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("director not found")))
		return
	}

	ratedBranches, err := server.ratedBranchesFunc(request.BrandId, monthPeriod, *saleTypes)
	var rate int32
	if ratedBranches != nil {
		for index, branch := range ratedBranches {
			if branch.BranchId == request.BranchId {
				rate = int32(index + 1.)
				break
			}
		}
	}

	ctx.JSON(http.StatusOK, entities2.BranchDashboardResponse{
		SalesStatisticsByTypes: saleStatisticByTypes,
		BestSalesManagers:      entities2.SalesManagerBranchItemsFromRatedUsers(bestSalesManagers),
		Branch: entities2.BranchModelResponse{
			Title:       branchInfo.Title,
			Description: branchInfo.Description,
		},
		GoalAchievementPercent: goalAchievementPercent * 100,
		Rating:                 rate,
		Profile: entities2.SimpleProfile{
			Id:       director[0].Id,
			Avatar:   director[0].AvatarPointer(),
			FullName: director[0].GetFullName(),
		},
	})
}

type GetBranchesRequest struct {
	BrandId int32 `json:"brand_id" form:"brand_id"`
}

type BranchesResp struct {
	Result []BranchResp `json:"result"`
}

func (server *Server) GetBranchesByBrand(ctx *gin.Context) {
	var request GetBranchesRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	branches, err := server.getBranchesByBrandFunc(request.BrandId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result := make([]BranchResp, 0)
	for _, branch := range branches {
		result = append(result, BranchResp{
			ID:    branch.BranchId,
			Title: branch.Title,
		})
	}

	ctx.JSON(http.StatusOK, BranchesResp{
		Result: result,
	})
}
