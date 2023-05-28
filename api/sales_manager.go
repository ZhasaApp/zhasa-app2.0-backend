package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	. "zhasa2.0/api/entities"
	"zhasa2.0/base"
	. "zhasa2.0/manager/entities"
	"zhasa2.0/manager/service"
	. "zhasa2.0/sale/entities"
	. "zhasa2.0/statistic/entities"
	"zhasa2.0/user/entities"
	tokenservice "zhasa2.0/user/service"
)

func getSalesManager(service tokenservice.TokenService, salesManagerService service.SalesManagerService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := tokenservice.Token(ctx.GetHeader("Authorization"))
		userData, err := service.VerifyToken(token)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		salesManager, err := salesManagerService.GetSalesManagerByUserId(userData.Id)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("sales manager not found"))
			return
		}

		log.Println(salesManager.Id)

		ctx.Set("sales_manager_id", int(salesManager.Id))
		ctx.Next()
	}
}

func (server *Server) createSalesManager(ctx *gin.Context) {
	var createUserBody CreateSalesManagerBody
	if err := ctx.ShouldBindJSON(&createUserBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	firstName, err := entities.NewName(createUserBody.FirstName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	lastName, err := entities.NewName(createUserBody.LastName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	phone, err := entities.NewPhone(createUserBody.Phone)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	createUserRequest := entities.CreateUserRequest{
		Phone:     *phone,
		FirstName: *firstName,
		LastName:  *lastName,
	}

	user, err := server.userService.GetUserByPhone(*phone)
	if user != nil && err == nil {
		err = server.salesManagerService.CreateSalesManager(user.Id, createUserBody.BranchId)

		if err != nil {
			fmt.Println(err)
			ctx.JSON(http.StatusBadRequest, errors.New("create sales manager error"))
			return
		}
		ctx.Status(http.StatusOK)
		return
	}

	err = server.userService.CreateUser(createUserRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err = server.userService.GetUserByPhone(*phone)
	if user == nil && err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.salesManagerService.CreateSalesManager(user.Id, createUserBody.BranchId)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, errors.New("create sales manager error"))
		return
	}
	ctx.Status(http.StatusOK)
}

func (server *Server) getSales(ctx *gin.Context) {
	var monthPagination MonthPaginationRequest
	if err := ctx.ShouldBindJSON(&monthPagination); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	salesManager, err := server.salesManagerService.GetSalesManagerByUserId(monthPagination.UserId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("sales manager not found")))
		return
	}

	salesCount, err := server.salesManagerService.GetSalesManagerSalesCount(salesManager.Id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("sales manager not found")))
		return
	}

	sales, err := server.salesManagerService.GetManagerSales(salesManager.Id, base.Pagination{
		PageSize: monthPagination.PageSize,
		Page:     monthPagination.Page,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	salesList := make([]SaleItemResponse, 0)
	for _, item := range *sales {
		salesList = append(salesList, SaleItemResponse{
			Id:     int32(item.Id),
			Title:  string(item.SaleDescription),
			Date:   item.SaleDate.Format("2006-01-02 15:04:05"),
			Amount: int64(item.SalesAmount),
			Type: SaleTypeResponse{
				Id:    int32(item.SaleType.Id),
				Title: item.SaleType.Title,
				Color: item.SaleType.Color,
			},
		})
	}

	hasNext := salesCount > monthPagination.PageSize*(monthPagination.Page+1)

	salesResponse := SalesResponse{
		Result:  salesList,
		Count:   int32(len(salesList)),
		HasNext: hasNext,
	}

	ctx.JSON(http.StatusOK, salesResponse)
}

func (server *Server) saveSale(ctx *gin.Context) {
	var saveSaleBody SaveSaleBody
	if err := ctx.ShouldBindJSON(&saveSaleBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	salesManagerId := ctx.GetInt("sales_manager_id")

	saleTypeId := saveSaleBody.SaleTypeId

	_, err := server.saleTypeService.GetSaleType(SaleTypeId(saleTypeId))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("sale type not found")))
		return
	}

	layout := "2006-01-02 15:04:05"
	parsedTime, err := time.Parse(layout, saveSaleBody.SaleDate)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	saleType, err := server.saleTypeService.GetSaleType(SaleTypeId(saveSaleBody.SaleTypeId))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("sale type not found")))
		return
	}

	sale := Sale{
		SaleManagerId:   SalesManagerId(salesManagerId),
		SaleType:        *saleType,
		SalesAmount:     SaleAmount(saveSaleBody.SaleAmount),
		SaleDate:        parsedTime,
		SaleDescription: SaleDescription(saveSaleBody.Description),
	}

	saleRes, err := server.salesManagerService.SaveSale(sale)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SaleItemResponse{
		Id:     int32(saleRes.Id),
		Title:  string(saleRes.SaleDescription),
		Date:   saleRes.SaleDate.Format("2006-01-02 15:04:05"),
		Amount: int64(saleRes.SalesAmount),
		Type: SaleTypeResponse{
			Id:    int32(saleRes.SaleType.Id),
			Title: saleRes.SaleType.Title,
			Color: saleRes.SaleType.Color,
		},
	})
}

func (server *Server) getSalesManagerDashboardStatistic(ctx *gin.Context) {
	var requestBody SalesManagerMonthStatisticRequestBody
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	salesManager, err := server.salesManagerService.GetSalesManagerByUserId(requestBody.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	period := MonthPeriod{
		MonthNumber: requestBody.Month,
		Year:        requestBody.Year,
	}
	fromDate, toDate := period.ConvertToTime()

	sums, err := server.salesManagerService.GetSalesManagerSums(fromDate, toDate, salesManager.Id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	dailyPeriod := DayPeriod{
		Day: time.Now(),
	}
	dayStart, dayEnd := dailyPeriod.ConvertToTime()

	dailySums, err := server.salesManagerService.GetSalesManagerSums(dayStart, dayEnd, salesManager.Id)

	totalDailySum := dailySums.TotalSum()
	totalPeriodSum := sums.TotalSum()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	goal, err := server.salesManagerService.GetSalesManagerGoal(fromDate, toDate, salesManager.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("no goal found")))
		return
	}
	totalPercent := base.NewPercent(int64(totalPeriodSum), int64(goal))
	dailyPercent := base.NewPercent(int64(totalDailySum), int64(goal))

	salesStatisticItemsByTypes := make([]SalesStatisticsByTypesItem, 0)

	for key, amount := range *sums {
		item := SalesStatisticsByTypesItem{
			Color:  key.Color,
			Title:  key.Title,
			Amount: int64(amount),
		}
		salesStatisticItemsByTypes = append(salesStatisticItemsByTypes, item)
	}

	sales, err := server.salesManagerService.GetManagerSales(salesManager.Id, base.Pagination{
		PageSize: 5,
		Page:     0,
	})

	salesResponse := make([]SaleItemResponse, 0)

	if err == nil {
		for _, item := range *sales {

			salesResponse = append(salesResponse, SaleItemResponse{
				Id:     int32(item.Id),
				Title:  string(item.SaleDescription),
				Date:   item.SaleDate.Format("2006-01-02 15:04:05"),
				Amount: int64(item.SalesAmount),
				Type: SaleTypeResponse{
					Id:    int32(item.SaleType.Id),
					Title: item.SaleType.Title,
					Color: item.SaleType.Color,
				},
			})
		}
	}

	dr := SalesManagerDashboardResponse{
		Profile: SalesManagerDashboardProfile{
			Avatar:   nil,
			FullName: salesManager.FirstName + " " + salesManager.LastName,
			Branch:   string(salesManager.Branch.Title),
		},
		OverallSalesStatistics: OverallSalesStatistic{
			Goal:     int64(goal),
			Achieved: int64(totalPeriodSum),
			Percent:  float64(totalPercent),
			GrowthPerDay: GrowthPerDay{
				Amount:  int64(totalDailySum),
				Percent: float64(dailyPercent),
			},
		},
		SalesStatisticsByTypes: salesStatisticItemsByTypes,
		LastSales:              salesResponse,
		Rating:                 int32(1),
	}
	ctx.JSON(http.StatusOK, dr)
}

func (server *Server) getYearStatistic(ctx *gin.Context) {
	var requestBody SalesManagerYearStatisticRequestBody
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	salesManager, err := server.salesManagerService.GetSalesManagerByUserId(requestBody.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	data, err := server.salesManagerService.GetSalesManagerYearMonthlyStatistic(salesManager.Id, requestBody.Year)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	response := make([]YearStatisticResponse, 0)
	for _, item := range *data {
		response = append(response, YearStatisticResponse{
			SaleType: SaleTypeResponse{
				Title: item.SaleType.Title,
				Color: "",
			},
			Month:  int32(item.Month),
			Amount: int64(item.Amount),
		})
	}

	ctx.JSON(http.StatusOK, response)
	return
}
