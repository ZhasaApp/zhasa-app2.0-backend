package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	entities2 "zhasa2.0/api/entities"
	"zhasa2.0/base"
	"zhasa2.0/rating"
	"zhasa2.0/statistic"
)

type GetSMDashboardRequest struct {
	UserId  int32 `json:"user_id" form:"user_id"`
	BrandId int32 `json:"brand_id" form:"brand_id"`
	Month   int32 `json:"month" form:"month"`
	Year    int32 `json:"year" form:"year"`
}

func (server *Server) SMDashboard(ctx *gin.Context) {
	var request GetSMDashboardRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	types, err := server.saleTypeRepo.GetSaleTypes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	monthPeriod := statistic.MonthPeriod{
		MonthNumber: request.Month,
		Year:        request.Year,
	}

	sTypeSums := make([]entities2.SalesStatisticsByTypesItem, 0)
	dashboardResponse := entities2.SalesManagerDashboardResponse{}

	var goalAchievementPercent float32
	ratioRows := make([]rating.RatioRow, 0)

	for _, saleType := range *types {
		amount, err := server.getSaleSumByUserBrandTypePeriodFunc(request.UserId, request.BrandId, saleType.Id, monthPeriod)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		goal := server.userBrandGoal(request.UserId, request.BrandId, saleType.Id, monthPeriod)

		sTypeSums = append(sTypeSums, entities2.SalesStatisticsByTypesItem{
			Color:    saleType.Color,
			Title:    saleType.Title,
			Achieved: amount,
			Goal:     goal,
		})

		ratioRows = append(ratioRows, rating.RatioRow{
			Achieved: amount,
			Goal:     goal,
			Gravity:  saleType.Gravity,
		})
	}

	goalAchievementPercent = rating.CalculateRatio(ratioRows)

	r, err := server.getUserRatingFunc(request.UserId, request.BrandId, statistic.MonthPeriod{
		MonthNumber: request.Month,
		Year:        request.Year,
	})

	dashboardResponse.Rating = r
	dashboardResponse.SalesStatisticsByTypes = sTypeSums
	dashboardResponse.GoalAchievementPercent = goalAchievementPercent * 100

	err = server.updateUserBrandRatio(request.UserId, request.BrandId, float64(goalAchievementPercent), monthPeriod)
	if err != nil {
		fmt.Println(err)
	}

	lastSales, err := server.salesByBrandUserFunc(request.UserId, request.BrandId, base.Pagination{
		PageSize: 5,
		Page:     0,
	})

	saleItems := entities2.SaleItemsFromSales(lastSales)

	dashboardResponse.LastSales = saleItems

	user, err := server.getUserByIdFunc(request.UserId)

	branch, err := server.getUserBranchFunc(request.UserId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	dashboardResponse.Profile = entities2.SalesManagerDashboardProfile{
		Id:       request.UserId,
		Avatar:   user.AvatarPointer(),
		FullName: user.GetFullName(),
		Branch: entities2.SimpleBranch{
			Id:    branch.ID,
			Title: branch.Title,
		},
	}

	ctx.JSON(http.StatusOK, dashboardResponse)
}
