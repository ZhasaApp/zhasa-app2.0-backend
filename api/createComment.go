package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type createCommentRequest struct {
	NewsId  int32  `json:"news_id" binding:"required"`
	Message string `json:"message" binding:"required"`
}

func (server *Server) CreateComment(ctx *gin.Context) {
	var req *createCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := int32(ctx.GetInt("user_id"))

	err := server.postRepository.CreateComment(userId, req.NewsId, req.Message)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}
