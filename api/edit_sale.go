package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"time"
	"zhasa2.0/api/entities"
	entities2 "zhasa2.0/statistic"
)

type EditSaleRequest struct {
	ID      int32  `json:"id"`
	Date    string `json:"date"`
	TypeID  int32  `json:"type_id"`
	Value   int64  `json:"value"`
	Title   string `json:"title"`
	BrandId int32  `json:"brand_id"`
}

func cleanDateTime(input string) string {
	// This regex captures date and time up to seconds, and discards everything after it.
	re := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{1,2}:\d{2}:\d{2})`)
	match := re.FindStringSubmatch(input)
	if len(match) > 1 {
		return match[1]
	}
	return input // return original if no match
}

func (server *Server) EditSale(ctx *gin.Context) {
	var requestBody EditSaleRequest
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := int32(ctx.GetInt("user_id"))

	layout := "2006-01-02 15:04:05"
	parsedTime, err := time.Parse(layout, cleanDateTime(requestBody.Date))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := server.saleEditFunc(requestBody.Value, parsedTime, requestBody.ID, userId, requestBody.TypeID, requestBody.Title)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	sType, err := server.saleTypeRepo.GetSaleType(requestBody.TypeID)
	if err != nil {
		fmt.Println(err)
		return
	}

	monthPeriod := entities2.MonthPeriod{
		MonthNumber: int32(parsedTime.Month()),
		Year:        int32(parsedTime.Year()),
	}

	goalAchievement, err := server.calculateUserBrandRatio(userId, requestBody.BrandId, monthPeriod)

	if err == nil {
		err := server.updateUserBrandRatio(userId, requestBody.BrandId, float64(goalAchievement), monthPeriod)
		if err != nil {
			fmt.Println(err)
		}
	}

	ctx.JSON(http.StatusOK, entities.SaleItemResponse{
		Id:     id,
		Title:  requestBody.Title,
		Date:   requestBody.Date,
		Amount: requestBody.Value,
		Type: entities.SaleTypeResponse{
			Id:        requestBody.TypeID,
			Title:     sType.Title,
			Color:     sType.Color,
			ValueType: sType.ValueType,
		},
	})
}
