package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/statistic/entities"
)

type DeleteSaleQuery struct {
	SaleId int32 `json:"id" form:"id"`
}

func (server *Server) DeleteSale(ctx *gin.Context) {
	var requestBody DeleteSaleQuery
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userId := int32(ctx.GetInt("user_id"))

	saleBrand, err := server.saleRepo.GetSaleBrandId(requestBody.SaleId)

	err = server.saleRepo.DeleteSale(requestBody.SaleId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.Status(http.StatusNoContent)

	ratio, err := server.calculateUserBrandRatio(userId, saleBrand.BrandID, entities.MonthPeriod{
		MonthNumber: int32(saleBrand.SaleDate.Month()),
		Year:        int32(saleBrand.SaleDate.Year()),
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	err = server.updateUserBrandRatio(userId, saleBrand.BrandID, float64(ratio), entities.MonthPeriod{
		MonthNumber: int32(saleBrand.SaleDate.Month()),
		Year:        int32(saleBrand.SaleDate.Year()),
	})
}
