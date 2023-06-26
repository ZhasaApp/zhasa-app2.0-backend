package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	salesStatisticItemsByTypes := make([]SalesStatisticsByTypesItem, 0)

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
			Avatar:      nil,
			FullName:    sm.FirstName + " " + sm.LastName,
			Ratio:       float64(sm.Ratio),
			BranchTitle: string(sm.Branch.Title),
			BranchId:    int32(sm.Branch.BranchId),
		})
	}

	dr := BranchDashboardResponse{
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
