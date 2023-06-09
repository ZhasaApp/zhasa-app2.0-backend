package api

import (
	"errors"
	"github.com/gin-gonic/gin"
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

	sm, err := server.salesManagerService.GetSalesManagerByUserId(userTokenData.Id)

	if sm != nil {
		response := UserProfileResponse{
			Id:       userTokenData.Id,
			Avatar:   nil,
			FullName: userTokenData.FirstName + " " + userTokenData.LastName,
			Phone:    userTokenData.Phone,
			Branch: BranchResponse{
				Id:          int32(sm.Branch.BranchId),
				Description: string(sm.Branch.Title),
			},
			Role: "sales_manager",
		}

		ctx.JSON(http.StatusOK, response)
		return
	}

	bd, err := server.directorService.GetBranchDirectorByUserId(UserId(userTokenData.Id))
	if bd != nil {
		response := UserProfileResponse{
			Id:       userTokenData.Id,
			Avatar:   nil,
			FullName: userTokenData.FirstName + " " + userTokenData.LastName,
			Phone:    userTokenData.Phone,
			Branch: BranchResponse{
				Id:          int32(bd.Branch.BranchId),
				Description: string(bd.Branch.Title),
			},
			Role: "branch_director",
		}

		ctx.JSON(http.StatusOK, response)
		return
	}
	response := UserProfileResponse{
		Id:       userTokenData.Id,
		Avatar:   nil,
		FullName: userTokenData.FirstName + " " + userTokenData.LastName,
		Phone:    userTokenData.Phone,
		Role:     "admin",
	}

	ctx.JSON(http.StatusOK, response)
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
		FirstName: string(user.FirstName),
		LastName:  string(user.LastName),
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
func (server *Server) createUser(ctx *gin.Context) {
	var createUserBody CreateUserBody
	if err := ctx.ShouldBindJSON(&createUserBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	firstName, err := NewName(createUserBody.FirstName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	lastName, err := NewName(createUserBody.LastName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	phone, err := NewPhone(createUserBody.Phone)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	request := CreateUserRequest{
		Phone:     *phone,
		FirstName: *firstName,
		LastName:  *lastName,
	}

	user, err := server.userService.GetUserByPhone(*phone)
	if user != nil && err == nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("user already exist")))
		return
	}

	err = server.userService.CreateUser(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}
