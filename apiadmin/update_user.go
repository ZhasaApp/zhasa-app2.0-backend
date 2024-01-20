package apiadmin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/user/entities"
)

type UpdateUserRequest struct {
	UserID    int32   `json:"user_id"`
	Phone     string  `json:"phone"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	RoleKey   string  `json:"role"`
	Brands    []int32 `json:"brand_ids"`
	BranchID  int32   `json:"branch_id"`
}

func (s *Server) UpdateUser(ctx *gin.Context) {
	var request UpdateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	firstName, err := entities.NewName(request.FirstName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	lastName, err := entities.NewName(request.LastName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	validatedPhone, err := entities.NewPhone(request.Phone)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err = s.updateUserFunc(request.UserID, *firstName, *lastName, *validatedPhone); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err = s.updateUserBrands(request.UserID, request.Brands); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err = s.updateUserBranchFunc(request.UserID, request.BranchID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err = s.updateUserRole(request.UserID, request.RoleKey); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
	})
}

type ChangeUsersRoleRequest struct {
	UserIds []int32 `json:"user_ids"`
	RoleKey string  `json:"new_role"`
}

func (s *Server) ChangeUsersRole(ctx *gin.Context) {
	var request ChangeUsersRoleRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	for _, userId := range request.UserIds {
		if err := s.updateUserRole(userId, request.RoleKey); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "users role changed successfully",
	})
}

type ChangeUsersBrandsRequest struct {
	UserIds   []int32 `json:"user_ids"`
	BrandsIds []int32 `json:"new_brand_ids"`
}

func (s *Server) ChangeUsersBrands(ctx *gin.Context) {
	var request ChangeUsersBrandsRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	for _, userId := range request.UserIds {
		if err := s.updateUserBrands(userId, request.BrandsIds); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "users brands changed successfully",
	})
}

type ChangeUsersBranchRequest struct {
	UserIds  []int32 `json:"user_ids"`
	BranchId int32   `json:"new_branch_id"`
}

func (s *Server) ChangeUsersBranch(ctx *gin.Context) {
	var request ChangeUsersBranchRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	for _, userId := range request.UserIds {
		if err := s.updateUserBranchFunc(userId, request.BranchId); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "users branch changed successfully",
	})
}
