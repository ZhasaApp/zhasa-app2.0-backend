package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/api/entities"
	"zhasa2.0/base"
)

type UserBrandYearStatisticRequestBody struct {
	UserId  int32 `form:"user_id" json:"user_id"`
	Year    int32 `form:"year" json:"year"`
	BrandId int32 `form:"brand_id" json:"brand_id"`
}

func (server *Server) GetUserBrandYearStatistic(ctx *gin.Context) {
	// retrieve year statistic for user with given request body
	var requestBody UserBrandYearStatisticRequestBody
	if err := ctx.ShouldBindQuery(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// use server.saleRepo to retrieve statistic from db
	stats, err := server.saleRepo.GetUserBrandMonthlyYearStatistic(requestBody.Year, requestBody.UserId, requestBody.BrandId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// convert stats to YearStatisticResultResponse and return
	response := entities.YearStatisticResultResponse{
		Result: make([]entities.YearStatisticResponse, 0),
	}
	for _, stat := range stats {
		response.Result = append(response.Result, entities.YearStatisticResponse{
			SaleType: entities.SaleTypeResponse{
				Id:        stat.SaleType.Id,
				Title:     stat.SaleType.Title,
				Color:     stat.SaleType.Color,
				ValueType: stat.SaleType.ValueType,
			},
			Month:  stat.Month,
			Amount: stat.Amount,
			Goal:   stat.Goal,
		})
	}
	ctx.JSON(http.StatusOK, response)
}

type GetUserSalesStatisticsRequest struct {
	UserId     int32 `form:"user_id" json:"user_id"`
	Year       int32 `form:"year" json:"year"`
	BrandId    int32 `form:"brand_id" json:"brand_id"`
	SaleTypeId int32 `form:"sale_type_id" json:"sale_type_id"`
}

func (server *Server) GetUserSalesStatistics(ctx *gin.Context) {
	var request GetUserSalesStatisticsRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	stats, err := server.saleRepo.GetUserBrandMonthlyYearStatistic(request.Year, request.UserId, request.BrandId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	result := make([]SalesStatisticsItem, 0)
	for _, stat := range stats {
		if stat.SaleType.Id != request.SaleTypeId {
			continue
		}
		rate := .0
		if stat.Goal != 0 {
			rate = float64(stat.Amount) / float64(stat.Goal) * 100.0
		}
		result = append(result, SalesStatisticsItem{
			Month:           stat.Month,
			ValueType:       stat.SaleType.ValueType,
			Achieved:        stat.Amount,
			Goal:            stat.Goal,
			GoalAchievement: BuildSuccessRateResp(rate),
		})
	}
	ctx.JSON(http.StatusOK, base.ArrayResponse[SalesStatisticsItem]{
		Result: result,
	})
}
