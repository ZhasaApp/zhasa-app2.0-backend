package apiadmin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type DeleteUserRequest struct {
	UserIds []int32 `json:"user_ids"`
}

func (s *Server) DeleteUsers(ctx *gin.Context) {
	var request DeleteUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	for _, userId := range request.UserIds {
		if err := s.addDisabledUserFunc(userId); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "users disabled successfully",
	})
}
