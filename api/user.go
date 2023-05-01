package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/user/entities"
	"zhasa2.0/user/service"
)

type CreateUserBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

type RequestAuthCodeBody struct {
	Phone string `json:"phone"`
}

type RequestAuthCodeResponse struct {
	OtpId int32 `json:"otp_id"`
}

type TryAuthBody struct {
	OtpId int32 `json:"otpId"`
	Otp   int32 `json:"otp"`
}

type AuthResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Token     string `json:"token"`
}

func (server *Server) tryAuth(ctx *gin.Context) {
	var request TryAuthBody
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	phone, err := entities.NewPhone(string(request.OtpId))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.authService.Login(*phone, entities.OtpCode(request.Otp))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userTokenData := service.UserTokenData{
		Id:        user.Id,
		FirstName: string(user.FirstName),
		LastName:  string(user.LastName),
		Email:     string(user.LastName),
	}
	token, err := server.tokenService.GenerateToken(&userTokenData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.New("error generating new token"))
	}

	response := AuthResponse{
		FirstName: userTokenData.FirstName,
		LastName:  userTokenData.LastName,
		Phone:     string(*phone),
		Token:     string(token),
	}
	ctx.JSON(http.StatusOK, response)
}

func (server *Server) requestAuthCode(ctx *gin.Context) {
	var request RequestAuthCodeBody
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	phone, err := entities.NewPhone(request.Phone)
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
		OtpId: id,
	}

	ctx.JSON(http.StatusOK, response)
}
func (server *Server) createUser(ctx *gin.Context) {
	var createUserBody CreateUserBody
	if err := ctx.ShouldBindJSON(&createUserBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	firstName, err := entities.NewName(createUserBody.FirstName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	lastName, err := entities.NewName(createUserBody.LastName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	phone, err := entities.NewPhone(createUserBody.Phone)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	request := entities.CreateUserRequest{
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
