package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/pkg/notify"
)

type createPostRequest struct {
	Title     string   `json:"title" binding:"required"`
	Body      string   `json:"body" binding:"required"`
	ImageUrls []string `json:"images"`
}

func (server *Server) CreatePost(ctx *gin.Context) {
	var req *createPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := ctx.GetInt("user_id")

	id, err := server.postRepository.CreatePost(req.Title, req.Body, int32(userId), req.ImageUrls)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	err = server.fbClient.SendPushNotification(ctx, notify.Message{
		Heading: "Опубликована новость",
		Message: req.Title,
		Payload: map[string]string{
			"deeplink": "doschamp://news?id=" + fmt.Sprintf("%d", id),
		},
	}, "news-topic")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send push notification",
		})
		return
	}

	ctx.Status(http.StatusNoContent)
	return
}
