package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/statistic"
)

type SetOwnerDashboardGoalRequest struct {
	BrandId    int32 `json:"brand_id"`
	Value      int64 `json:"value"`
	Month      int32 `json:"month"`
	Year       int32 `json:"year"`
	SaleTypeID int32 `json:"sale_type_id"`
}

func (server *Server) SetOwnerDashboardGoal(ctx *gin.Context) {
	var request SetOwnerDashboardGoalRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	period := statistic.MonthPeriod{
		MonthNumber: request.Month,
		Year:        request.Year,
	}

	err := server.setBrandSaleTypeGoal(request.BrandId, request.SaleTypeID, request.Value, period)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

type GetOwnerDashboardBySaleTypesRequest struct {
	Month   int32 `json:"month"`
	Year    int32 `json:"year"`
	BrandId int32 `json:"brand_id"`
}

type SaleType struct {
	Title     string `json:"title"`
	Color     string `json:"color"`
	ValueType string `json:"value_type"`
}

type OwnerDashboardBySaleTypesItem struct {
	SaleType SaleType `json:"sale_type"`
	Achieved int64    `json:"achieved"`
	Goal     int64    `json:"goal"`
}

func (server *Server) GetOwnerDashboardBySaleTypes(ctx *gin.Context) {
	var request GetOwnerDashboardBySaleTypesRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	period := statistic.MonthPeriod{
		MonthNumber: request.Month,
		Year:        request.Year,
	}

	saleTypes, err := server.saleTypeRepo.GetSaleTypes()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result := make([]OwnerDashboardBySaleTypesItem, 0)

	for _, saleType := range *saleTypes {
		achieved, _ := server.getBrandSaleSumFunc(request.BrandId, saleType.Id, period)
		goal, _ := server.getBrandGoalFunc(request.BrandId, saleType.Id, period)

		result = append(result, OwnerDashboardBySaleTypesItem{
			SaleType: SaleType{
				Title:     saleType.Title,
				Color:     saleType.Color,
				ValueType: saleType.ValueType,
			},
			Achieved: achieved,
			Goal:     goal,
		})
	}

	ctx.JSON(http.StatusOK, result)
}
