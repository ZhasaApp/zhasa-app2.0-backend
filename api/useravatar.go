package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strings"
	. "zhasa2.0/user/entities"
)

type UpdateUserProfileAvatarRequest struct {
	Url string `json:"url"`
}

func isImageUrl(url string) bool {
	extension := strings.ToLower(filepath.Ext(url))
	switch extension {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		return true
	}
	return false
}

func (server *Server) UploadUserAvatar(ctx *gin.Context) {
	var request UpdateUserProfileAvatarRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !isImageUrl(request.Url) {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("not an image url")))
		return
	}

	userId := ctx.GetInt("user_id")

	if userId == 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("user not found")))
		return
	}
	err := server.userService.UploadAvatar(UserId(userId), request.Url)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (server Server) DeleteAvatar(ctx *gin.Context) {
	userId := ctx.GetInt("user_id")

	if userId == 0 {
		ctx.JSON(http.StatusUnauthorized, errors.New("user not found"))
		return
	}

	err := server.userService.DeleteAvatar(UserId(userId))

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
