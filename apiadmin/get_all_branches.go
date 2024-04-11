package apiadmin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	generated "zhasa2.0/db/sqlc"
)

type BranchItem struct {
	Id          int32  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GetAllBranchesRequest struct {
	Search    string `json:"search" form:"search"`
	SortType  string `json:"sort_type" form:"sort_type"`
	SortField string `json:"sort_field" form:"sort_field"`
}

type BranchesResponse struct {
	Result []BranchItem `json:"result"`
}

func (s *Server) GetAllBranches(ctx *gin.Context) {
	var req GetAllUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	branches, err := s.getBranchesFiltered(
		generated.GetBranchesSearchParams{
			Search:    req.Search,
			SortType:  req.SortType,
			SortField: req.SortField,
		})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var response BranchesResponse
	for _, branch := range branches {
		response.Result = append(response.Result, BranchItem{
			Id:          branch.BranchId,
			Title:       branch.Title,
			Description: branch.Description,
		})
	}
	ctx.JSON(http.StatusOK, response)
}
