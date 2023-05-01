package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/branch/entities"
)

type createBranchBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (server *Server) createBranch(ctx *gin.Context) {
	var createBranchBody createBranchBody
	if err := ctx.ShouldBindJSON(&createBranchBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	title := entities.NewBranchTitle(createBranchBody.Title)
	description := entities.NewBranchDescription(createBranchBody.Description)

	request := entities.CreateBranchRequest{
		Title:       title,
		Description: description,
	}
	err := server.branchService.CreateBranch(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}
