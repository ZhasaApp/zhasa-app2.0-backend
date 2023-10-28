package apiadmin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BranchItem struct {
	Id          int32  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type BranchesResponse struct {
	Result []BranchItem `json:"result"`
}

func (s *Server) GetAllBranches(ctx *gin.Context) {
	branches, err := s.getAllBranchesFunc()
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
