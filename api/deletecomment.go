package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type deleteCommentRequest struct {
	CommentId int32 `json:"id" binding:"required"`
}

func (server Server) DeleteComment(ctx *gin.Context) {
	var req *deleteCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.postRepository.DeleteComment(req.CommentId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}
