package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sort"
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

func (server *Server) getBranchDashboardStatistic(ctx *gin.Context) {
	var requestBody BranchMonthStatisticRequestBody
	if err := ctx.ShouldBindQuery(&requestBody); err != nil {
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

	salesManagers, err := server.branchService.GetBranchRankedSalesManagers(period, branch.BranchId, Pagination{
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
			Avatar:      sm.GetAvatarPointer(),
			FullName:    sm.FirstName + " " + sm.LastName,
			Ratio:       float64(sm.Ratio),
			BranchTitle: string(sm.Branch.Title),
			BranchId:    int32(sm.Branch.BranchId),
		})
	}

	data, err := server.branchService.GetBranchSalesSums(fromDate, toDate, branch.BranchId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for st, sum := range *data {
		goal, err := server.branchService.GetBranchGoal(period, branch.BranchId, st.Id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		salesStatisticItemsByTypes = append(salesStatisticItemsByTypes, SalesStatisticsByTypesItem{
			Color:    st.Color,
			Title:    st.Title,
			Achieved: int64(sum),
			Goal:     int64(goal),
		})
	}

	director, err := server.directorService.GetBranchDirectorByBranchId(branch.BranchId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	branches, err := server.branchService.GetBranches(period)
	var branchGoalAchievement Percent
	rating := 1
	if err != nil {
		fmt.Println(err)
	}
	sort.Slice(branches, func(i, j int) bool {
		return branches[i].GoalAchievement > branches[j].GoalAchievement
	})

	if branches != nil {
		for index, item := range branches {
			if branch.BranchId == item.BranchId {
				branchGoalAchievement = item.GoalAchievement
				rating = index + 1
			}
		}
	}

	dr := BranchDashboardResponse{
		SalesStatisticsByTypes: salesStatisticItemsByTypes,
		BestSalesManagers:      bestSalesManagers,
		Branch: BranchModelResponse{
			Title:       string(branch.Title),
			Description: string(branch.Description),
		},
		Profile: SimpleProfile{
			Avatar:   director.AvatarPointer(),
			FullName: director.GetFullName(),
			Id:       director.Id,
		},
		GoalAchievementPercent: float32(branchGoalAchievement),
		Rating:                 int32(rating),
	}
	ctx.JSON(http.StatusOK, dr)
}

func (server *Server) getBranchYearStatistic(ctx *gin.Context) {
	var requestBody BranchYearStatisticRequestBody
	if err := ctx.ShouldBindQuery(&requestBody); err != nil {
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
		goal, err := server.branchService.GetBranchGoal(MonthPeriod{
			MonthNumber: int32(item.Month),
			Year:        requestBody.Year,
		}, BranchId(requestBody.BranchId), item.SaleType.Id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		response = append(response, YearStatisticResponse{
			SaleType: SaleTypeResponse{
				Title:     item.SaleType.Title,
				Color:     item.SaleType.Color,
				Id:        int32(item.SaleType.Id),
				ValueType: item.SaleType.ValueType,
			},
			Month:  int32(item.Month),
			Amount: int64(item.Amount),
			Goal:   int64(goal),
		})
	}

	ctx.JSON(http.StatusOK, YearStatisticResultResponse{
		Result: response,
	})
}
