package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ToggleLikeRequest struct {
	Id int32 `json:"id"`
}

func (server *Server) ToggleLike(ctx *gin.Context) {
	var request ToggleLikeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := int32(ctx.GetInt("user_id"))

	isLiked, err := server.postRepository.IsUserLikedPost(userId, request.Id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if isLiked {
		err := server.postRepository.DeleteLike(userId, request.Id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.Status(http.StatusNoContent)
		return
	} else {
		err := server.postRepository.AddLike(userId, request.Id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.Status(http.StatusNoContent)
		return
	}
}
