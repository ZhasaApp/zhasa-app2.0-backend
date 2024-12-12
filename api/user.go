package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	. "zhasa2.0/api/entities"
	. "zhasa2.0/user/entities"
	"zhasa2.0/user/service"
)

func verifyToken(tokenService service.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := service.Token(ctx.GetHeader("Authorization"))
		userData, err := tokenService.VerifyToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid token")))
			ctx.Abort()
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

	user, err := server.getUserByIdFunc(userTokenData.Id)
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
		Id:            user.Id,
		Avatar:        user.AvatarPointer(),
		FullName:      user.GetFullName(),
		Phone:         string(user.Phone),
		About:         user.About,
		Branch:        branchResponse,
		Role:          user.UserRole.Key,
		Branches:      nil,
		WorkStartDate: user.CreatedAt.Format("2006-01-02 15:04:05"),
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

type SearchUsersRequest struct {
	Search string `json:"search" form:"search"`
}

type SearchUserItem struct {
	Id       int32   `json:"id"`
	Avatar   *string `json:"avatar"`
	FullName string  `json:"full_name"`
}

type SearchUsersResponse struct {
	Result []SearchUserItem `json:"result"`
}

func (server *Server) SearchUsers(ctx *gin.Context) {
	var request SearchUsersRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	users, err := server.searchUsersFunc(request.Search)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := SearchUsersResponse{
		Result: make([]SearchUserItem, 0),
	}
	for _, user := range users {
		response.Result = append(response.Result, SearchUserItem{
			Id:       user.Id,
			Avatar:   user.AvatarPointer(),
			FullName: user.GetFullName(),
		})
	}

	ctx.JSON(http.StatusOK, response)
}

type GetUserResponse struct {
	Id            int32           `json:"id"`
	Avatar        *string         `json:"avatar"`
	FullName      string          `json:"full_name"`
	Branch        *BranchResponse `json:"branch"`
	About         string          `json:"about"`
	Role          string          `json:"role"`
	WorkStartDate string          `json:"work_start_date"`
}

func (server *Server) GetUser(ctx *gin.Context) {
	idRaw := ctx.Query("id")

	id, err := strconv.ParseInt(idRaw, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid id")))
		return
	}

	user, err := server.getUserByIdFunc(int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	branch, err := server.getUserBranchFunc(user.Id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

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

	response := GetUserResponse{
		Id:            user.Id,
		Avatar:        user.AvatarPointer(),
		FullName:      user.GetFullName(),
		Branch:        branchResponse,
		Role:          user.UserRole.Key,
		WorkStartDate: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if user.About != nil {
		response.About = *user.About
	}

	ctx.JSON(http.StatusOK, response)
}

type UpdateUserAboutRequest struct {
	About *string `json:"about"`
}

func (server *Server) UpdateUserAbout(ctx *gin.Context) {
	userIDRaw, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("user not found")))
		return
	}

	userID := userIDRaw.(int)

	var request UpdateUserAboutRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.updateUserProfileAbout(int32(userID), request.About)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
