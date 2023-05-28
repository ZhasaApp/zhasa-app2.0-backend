package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	. "zhasa2.0/api/entities"
	. "zhasa2.0/base"
	. "zhasa2.0/branch/entities"
	. "zhasa2.0/statistic/entities"
)

type createBranchBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Key         string `json:"key"`
}

func (server *Server) createBranch(ctx *gin.Context) {
	var createBranchBody createBranchBody
	if err := ctx.ShouldBindJSON(&createBranchBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	title := NewBranchTitle(createBranchBody.Title)
	description := NewBranchDescription(createBranchBody.Description)
	key := NewBranchKey(createBranchBody.Key)
	request := CreateBranchRequest{
		Title:       title,
		Description: description,
		Key:         key,
	}
	err := server.branchService.CreateBranch(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}

func (server *Server) getBranches(ctx *gin.Context) {
	branches, err := server.branchService.GetBranches()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, branches)
}

func (server *Server) getBranchDashboardStatistic(ctx *gin.Context) {
	var requestBody BranchMonthStatisticRequestBody
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	branch, err := server.branchService.GetBranchById(BranchId(requestBody.BranchId))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	period := MonthPeriod{
		MonthNumber: requestBody.Month,
		Year:        requestBody.Year,
	}
	fromDate, toDate := period.ConvertToTime()

	sums, err := server.branchService.GetBranchSalesSums(fromDate, toDate, branch.BranchId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	dailyPeriod := DayPeriod{
		Day: time.Now(),
	}
	dayStart, dayEnd := dailyPeriod.ConvertToTime()

	dailySums, err := server.branchService.GetBranchSalesSums(dayStart, dayEnd, branch.BranchId)

	totalDailySum := dailySums.TotalSum()
	totalPeriodSum := sums.TotalSum()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	goal, err := server.branchService.GetBranchGoal(fromDate, toDate, branch.BranchId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("no goal found")))
		return
	}
	totalPercent := NewPercent(int64(totalPeriodSum), int64(goal))
	dailyPercent := NewPercent(int64(totalDailySum), int64(goal))

	salesStatisticItemsByTypes := make([]SalesStatisticsByTypesItem, 0)

	for key, amount := range *sums {
		item := SalesStatisticsByTypesItem{
			Color:  "",
			Title:  key.Title,
			Amount: int64(amount),
		}
		salesStatisticItemsByTypes = append(salesStatisticItemsByTypes, item)
	}

	bestSalesManagers := make([]SalesManagerBranchItem, 0)

	salesManagers, err := server.branchService.GetBranchRankedSalesManagers(fromDate, toDate, branch.BranchId, Pagination{
		PageSize: 3,
		Page:     0,
	})

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for _, sm := range *salesManagers {
		bestSalesManagers = append(bestSalesManagers, SalesManagerBranchItem{
			Id:          int32(sm.UserId),
			Avatar:      sm.AvatarUrl,
			FullName:    sm.FirstName + " " + sm.LastName,
			Ratio:       float64(sm.Ratio),
			BranchTitle: string(sm.Branch.Title),
			BranchId:    int32(sm.Branch.BranchId),
		})
	}

	dr := BranchDashboardResponse{
		OverallSaleStatistics: OverallSalesStatistic{
			Goal:     int64(goal),
			Achieved: int64(totalPeriodSum),
			Percent:  float64(totalPercent),
			GrowthPerDay: GrowthPerDay{
				Amount:  int64(totalDailySum),
				Percent: float64(dailyPercent),
			},
		},
		SaleStatisticsByTypes: salesStatisticItemsByTypes,
		BestSalesManagers:     bestSalesManagers,
	}
	ctx.JSON(http.StatusOK, dr)
}

func (server *Server) getBranchYearStatistic(ctx *gin.Context) {
	var requestBody BranchYearStatisticRequestBody
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	data, err := server.branchService.GetBranchYearStatistic(BranchId(requestBody.BranchId), requestBody.Year)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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
}
