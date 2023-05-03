package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/branch/entities"
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

	title := entities.NewBranchTitle(createBranchBody.Title)
	description := entities.NewBranchDescription(createBranchBody.Description)
	key := entities.NewBranchKey(createBranchBody.Key)
	request := entities.CreateBranchRequest{
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
