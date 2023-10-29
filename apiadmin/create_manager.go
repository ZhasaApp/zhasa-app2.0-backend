package apiadmin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateManagerRequest struct {
	UserId   int32   `json:"user_id"`
	BranchId int32   `json:"branch_id"`
	Brands   []int32 `json:"brands"`
}

func (s *Server) CreateManager(ctx *gin.Context) {
	var request CreateManagerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.makeUserAsManagerFunc(request.UserId, request.BranchId, request.Brands)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}
