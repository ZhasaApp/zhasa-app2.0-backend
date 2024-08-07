package apiadmin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateGoodRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *Server) CreateGood(ctx *gin.Context) {
	var request CreateGoodRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := s.createGoodFunc(request.Name, request.Description)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}
