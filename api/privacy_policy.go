package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) PrivacyPolicy(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "policy.html", nil)
}
