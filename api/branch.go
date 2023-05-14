package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	. "zhasa2.0/api/entities"
	. "zhasa2.0/branch/entities"
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

func (server *Server) GetBranches(ctx *gin.Context) {
	branches, err := server.branchService.GetBranches()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, branches)
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
