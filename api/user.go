package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	. "zhasa2.0/api/entities"
	. "zhasa2.0/user/entities"
	"zhasa2.0/user/service"
)

func verifyToken(tokenService service.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := service.Token(ctx.GetHeader("Authorization"))
		userData, err := tokenService.VerifyToken(token)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		ctx.Set("user_id", int(userData.Id))
		ctx.Next()
	}
}

func (server *Server) getUserProfile(ctx *gin.Context) {
	token := service.Token(ctx.GetHeader("Authorization"))
	userTokenData, err := server.tokenService.VerifyToken(token)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("invalid token"))
		return
	}

	user, err := server.userRepo.GetUserById(userTokenData.Id)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("user not found")))
		return
	}

	branch, err := server.getUserBranchFunc(user.Id)

	var branchResponse *BranchResponse

	if branch != nil {
		brands, err := server.getBranchBrands(branch.ID)
		if err != nil {
			fmt.Println(err)
			log.Fatal("no brands for branch")
		}
		branchResponse = &BranchResponse{
			Id:          branch.ID,
			Description: branch.Title,
			Brands:      BrandItemsFromBrands(brands),
		}
	}

	ctx.JSON(http.StatusOK, UserProfileResponse{
		Id:       user.Id,
		Avatar:   user.AvatarPointer(),
		FullName: user.GetFullName(),
		Phone:    string(user.Phone),
		Branch:   branchResponse,
		Role:     user.UserRole.Key,
		Branches: nil,
	})
}

func (server *Server) tryAuth(ctx *gin.Context) {
	var request TryAuthBody
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.authService.Login(OtpId(request.OtpId), OtpCode(request.Otp))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userTokenData := service.UserTokenData{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     string(user.Phone),
	}
	token, err := server.tokenService.GenerateToken(&userTokenData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.New("error generating new token"))
	}

	authResponse := AuthResponse{
		Token: string(token),
	}
	ctx.JSON(http.StatusOK, authResponse)
}

func (server *Server) requestAuthCode(ctx *gin.Context) {
	var request RequestAuthCodeBody
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	phone, err := NewPhone(request.Phone)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := server.authService.RequestCode(*phone)

	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	response := RequestAuthCodeResponse{
		OtpId: int32(id),
	}

	ctx.JSON(http.StatusOK, response)
}
