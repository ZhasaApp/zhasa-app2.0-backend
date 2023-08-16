package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	service2 "zhasa2.0/branch_director/service"
	"zhasa2.0/user/entities"
	token_service "zhasa2.0/user/service"
)

func getBranchDirector(service token_service.TokenService, branchDirectorService service2.BranchDirectorService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := token_service.Token(ctx.GetHeader("Authorization"))
		userData, err := service.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("invalid token")))
			return
		}

		director, err := branchDirectorService.GetBranchDirectorByUserId(entities.UserId(userData.Id))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		log.Println(director[0].Id)

		ctx.Set("director_id", int(director[0].Id))
		ctx.Next()
	}
}
