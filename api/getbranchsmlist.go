package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type getBranchUsersRequest struct {
	BranchID int32 `json:"branch_id" form:"branch_id"`
	BrandID  int32 `json:"brand_id" form:"brand_id"`
	Month    int32 `json:"month" form:"month"`
	Year     int32 `json:"year" form:"year"`
}

type branchUser struct {
	ID       int32   `json:"id"`
	Avatar   *string `json:"avatar,omitempty"`
	FullName string  `json:"full_name"`
}

type branchUsersResponse struct {
	Result []branchUser `json:"result"`
}

func (server *Server) GetBranchSalesManagerList(ctx *gin.Context) {
	var request getBranchUsersRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	smList, err := server.getUsersByBranchBrandRoleFunc(request.BranchID, request.BrandID, 2)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	branchUsers := make([]branchUser, 0)

	for _, sm := range smList {
		branchUsers = append(branchUsers, branchUser{
			ID:       sm.Id,
			Avatar:   sm.AvatarPointer(),
			FullName: sm.GetFullName(),
		})
	}

	ctx.JSON(http.StatusOK, branchUsersResponse{
		Result: branchUsers,
	})
}
