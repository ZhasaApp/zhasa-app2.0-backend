package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	entities2 "zhasa2.0/api/entities"
	"zhasa2.0/api/rating"
	"zhasa2.0/date"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/statistic/entities"
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

	userBrand, err := server.getUserBrandFunc(request.UserId, request.BrandId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("user brand not found")))
		return
	}

	monthPeriod := entities.MonthPeriod{
		MonthNumber: request.Month,
		Year:        request.Year,
	}

	from, to := monthPeriod.ConvertToTime()
	sTypeSums := make([]entities2.SalesStatisticsByTypesItem, 0)
	dashboardResponse := entities2.SalesManagerDashboardResponse{}

	var goalAchievementPercent float32
	ratioRows := make([]rating.RatioRow, 0)

	for _, saleType := range *types {
		amount, err := server.saleRepo.GetSumByUserIdBrandIdPeriodSaleTypeId(generated.GetSaleSumByUserIdBrandIdPeriodSaleTypeIdParams{
			ID:         request.UserId,
			BrandID:    request.BrandId,
			SaleDate:   from,
			SaleDate_2: to,
			SaleTypeID: saleType.Id,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		goal := server.userBrandGoal(generated.GetUserBrandGoalParams{
			UserBrand:  userBrand,
			SaleTypeID: saleType.Id,
			FromDate:   from,
			FromDate_2: to,
		})

		sTypeSums = append(sTypeSums, entities2.SalesStatisticsByTypesItem{
			Color:    saleType.Color,
			Title:    saleType.Title,
			Achieved: amount,
			Goal:     goal,
		})

		ratioRows = append(ratioRows, rating.RatioRow{
			Amount:  amount,
			Goal:    goal,
			Gravity: saleType.Gravity,
		})
	}

	goalAchievementPercent = rating.CalculateRatio(ratioRows)

	r, err := server.getUserRatingFunc(request.UserId, request.BrandId, entities.MonthPeriod{
		MonthNumber: request.Month,
		Year:        request.Year,
	})

	dashboardResponse.Rating = r
	dashboardResponse.SalesStatisticsByTypes = sTypeSums
	dashboardResponse.GoalAchievementPercent = goalAchievementPercent

	lastSales, err := server.saleRepo.GetSalesByBrandIdAndUserId(generated.GetSalesByBrandIdAndUserIdParams{
		BrandID: request.BrandId,
		UserID:  request.UserId,
		Limit:   5,
		Offset:  0,
	})

	saleItems := make([]entities2.SaleItemResponse, 0)
	for _, sale := range lastSales {
		saleItems = append(saleItems, entities2.SaleItemResponse{
			Id:     0,
			Title:  sale.Description,
			Date:   date.ConvertTimeToStringISO(sale.SaleDate),
			Amount: 0,
			Type: entities2.SaleTypeResponse{
				Id:        sale.SaleTypeID,
				Title:     sale.SaleTypeTitle,
				Color:     sale.Color,
				ValueType: string(sale.ValueType),
			},
		})
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
}
