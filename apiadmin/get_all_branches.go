package apiadmin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/branch/entities"
	generated "zhasa2.0/db/sqlc"
)

type BranchItem struct {
	Id          int32  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Brands      string `json:"brands"`
}

type GetAllBranchesRequest struct {
	Search    string `json:"search" form:"search"`
	SortType  string `json:"sort_type" form:"sort_type"`
	SortField string `json:"sort_field" form:"sort_field"`
}

type BranchesResponse struct {
	Result  []BranchItem `json:"result"`
	HasNext bool         `json:"has_next"`
	Count   int64        `json:"count"`
}

func (s *Server) GetAllBranches(ctx *gin.Context) {
	var req GetAllUsersRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.PageSize == 0 {
		req.PageSize = 10
	}

	if req.Page < 0 {
		req.Page = 0
	}

	var branches []entities.Branch

	switch req.SortType {
	case "asc":
		branches, err = s.getBranchesFilteredAsc(generated.GetBranchesSearchAscParams{
			Search: req.Search,
			Limit:  req.PageSize,
			Offset: req.Page,
		})
	default:
		branches, err = s.getBranchesFilteredDesc(generated.GetBranchesSearchDescParams{
			Search: req.Search,
			Limit:  req.PageSize,
			Offset: req.Page,
		})
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	count, err := s.getBranchesFilteredCount(req.Search)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	response := BranchesResponse{
		Result:  make([]BranchItem, 0, len(branches)),
		HasNext: int64(req.Page*req.PageSize) < count,
		Count:   count,
	}

	for _, branch := range branches {
		response.Result = append(response.Result, BranchItem{
			Id:          branch.BranchId,
			Title:       branch.Title,
			Description: branch.Description,
			Brands:      branch.Brands,
		})
	}
	ctx.JSON(http.StatusOK, response)
}
