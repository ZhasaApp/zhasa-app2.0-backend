package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	. "zhasa2.0/api/entities"
	. "zhasa2.0/statistic/entities"
)

func (server *Server) GetBranchList(ctx *gin.Context) {
	var monthPagination MonthPaginationRequest
	if err := ctx.ShouldBindQuery(&monthPagination); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	//
	//period := MonthPeriod{
	//	MonthNumber: monthPagination.Month,
	//	Year:        monthPagination.Year,
	//}
	//

}
