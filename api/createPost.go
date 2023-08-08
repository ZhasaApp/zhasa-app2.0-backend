package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type createPostRequest struct {
	Title     string   `json:"title" binding:"required"`
	Body      string   `json:"body" binding:"required"`
	ImageUrls []string `json:"image_urls"`
}

func (server Server) CreatePost(ctx *gin.Context) {
	var req *createPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := ctx.GetInt("user_id")

	err := server.postRepository.CreatePost(req.Title, req.Body, int32(userId), req.ImageUrls)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.Status(http.StatusNoContent)
	return
}
