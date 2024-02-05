package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (server *Server) signup(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}
