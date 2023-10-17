package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/rating"
	"zhasa2.0/statistic"
)

type GetRatedBranches struct {
	BrandId int32 `json:"brand_id" form:"brand_id"`
	Month   int32 `json:"month" form:"month"`
	Year    int32 `json:"year" form:"year"`
}

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

func (server *Server) GetRatedBranches(ctx *gin.Context) {
	var request GetRatedBranches
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	period := statistic.MonthPeriod{
		MonthNumber: request.Month,
		Year:        request.Year,
	}

	brandBranches, err := server.getBranchesByBrandFunc(request.BrandId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result := make([]BranchRatingItem, 0)

	sTypes, err := server.saleTypeRepo.GetSaleTypes()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	for _, branch := range brandBranches {
		branchRatioRows := make([]rating.RatioRow, 0)
		for _, saleType := range *sTypes {
			salesSum, _ := server.getBranchBrandSaleSumFunc(branch.BranchId, request.BrandId, saleType.Id, period)

			goal, _ := server.getBranchBrandGoalFunc(branch.BranchId, saleType.Id, period)

			branchRatioRows = append(branchRatioRows, rating.RatioRow{
				Achieved: salesSum,
				Goal:     goal,
				Gravity:  saleType.Gravity,
			})
		}
		result = append(result, BranchRatingItem{
			ID:                     branch.BranchId,
			Title:                  branch.Title,
			Description:            branch.Description,
			GoalAchievementPercent: rating.CalculateRatio(branchRatioRows) * 100,
		})
	}

	ctx.JSON(http.StatusOK, BranchesResponse{
		Result: result,
	})
}
