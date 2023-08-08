package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type deletePostRequest struct {
	PostId int32 `json:"id" binding:"required"`
}

func (server Server) DeletePost(ctx *gin.Context) {
	var req *deletePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.postRepository.DeletePost(req.PostId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}
