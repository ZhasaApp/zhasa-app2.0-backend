package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	entities2 "zhasa2.0/branch_director/entities"
	service2 "zhasa2.0/branch_director/service"
	entities3 "zhasa2.0/manager/entities"
	entities4 "zhasa2.0/sale/entities"
	entities5 "zhasa2.0/statistic/entities"
	"zhasa2.0/user/entities"
	token_service "zhasa2.0/user/service"
)

func getBranchDirector(service token_service.TokenService, branchDirectorService service2.BranchDirectorService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := token_service.Token(ctx.GetHeader("Authorization"))
		userData, err := service.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("invalid token")))
			return
		}

		salesManager, err := branchDirectorService.GetBranchDirectorByUserId(entities.UserId(userData.Id))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		log.Println(salesManager.Id)

		ctx.Set("sales_manager_id", int(salesManager.Id))
		ctx.Next()
	}
}

type createSalesManagerGoalBody struct {
	SalesManagerId int32 `json:"sales_manager_id"`
	Amount         int64 `json:"amount"`
	Year           int32 `json:"year"`
	Month          int32 `json:"month"`
}

func (server *Server) createSaleGoalForSalesManager(ctx *gin.Context) {
	var body createSalesManagerGoalBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if body.Month < 1 || body.Month > 12 {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid month number")))
		return
	}

	period := entities5.MonthPeriod{
		MonthNumber: body.Month,
		Year:        body.Year,
	}
	fromDate, toDate := period.ConvertToTime()

	goal := entities2.SalesManagerGoal{
		SalesManagerId: entities3.SalesManagerId(body.SalesManagerId),
		FromDate:       fromDate,
		ToDate:         toDate,
		Amount:         entities4.SaleAmount(body.Amount),
	}
	err := server.directorService.CreateSalesManagerGoal(goal)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}
