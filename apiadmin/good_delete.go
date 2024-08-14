package apiadmin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type DeleteGoodRequest struct {
	GoodId int32 `json:"good_id" form:"good_id" binding:"required"`
}

func (server *Server) DeleteGood(ctx *gin.Context) {
	var deleteGoodRequest DeleteGoodRequest
	if err := ctx.ShouldBindQuery(&deleteGoodRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.deleteGoodFunc(deleteGoodRequest.GoodId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}
