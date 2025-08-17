package api

import "github.com/gin-gonic/gin"

type AddUserTokenRequest struct {
	Token string `json:"token"`
}

func (s *Server) AddUserToken(ctx *gin.Context) {
	var request AddUserTokenRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	userId := ctx.GetInt("user_id")
	if userId == 0 {
		ctx.JSON(401, gin.H{"error": "User not found"})
		return
	}

	err := s.fbClient.SubscribeToTopic(ctx, request.Token, "news-topic")
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to subscribe to topic"})
		return
	}

	err = s.addUserTokenFunc(int32(userId), request.Token)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to add user token"})
		return
	}

	ctx.Status(200)
}
