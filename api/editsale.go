package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	. "zhasa2.0/manager/entities"
	. "zhasa2.0/statistic/entities"
)

type EditSaleRequest struct {
	ID     int32     `json:"id"`
	Date   time.Time `json:"date"`
	TypeID int32     `json:"type_id"`
	Value  int64     `json:"value"`
	Title  string    `json:"title"`
}

func (server Server) EditSale(ctx *gin.Context) {
	var requestBody EditSaleRequest
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	deletedSale, err := server.salesManagerService.EditSale(EditSaleBody{
		ID:     requestBody.ID,
		Date:   requestBody.Date,
		TypeID: requestBody.TypeID,
		Value:  requestBody.Value,
		Title:  requestBody.Title,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	server.salesManagerService.UpdateRatio(deletedSale.SaleManagerId, MonthPeriod{
		MonthNumber: int32(deletedSale.SaleDate.Month()),
		Year:        int32(deletedSale.SaleDate.Year()),
	})

	ctx.Status(http.StatusNoContent)
}
