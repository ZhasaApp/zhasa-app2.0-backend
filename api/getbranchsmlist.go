package api

import (
	"github.com/gin-gonic/gin"
)

type GetBranchUsersRequest struct {
	BranchID int32 `json:"branch_id" form:"branch_id"`
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
	//var request GetBranchUsersRequest
	//if err := ctx.ShouldBindQuery(&request); err != nil {
	//	ctx.JSON(http.StatusBadRequest, errorResponse(err))
	//	return
	//}
	//
	//smList, err := server.branchService.GetBranchRankedSalesManagers(MonthPeriod{
	//	MonthNumber: request.Month,
	//	Year:        request.Year,
	//}, entities.BranchId(request.BranchID), Pagination{
	//	PageSize: 20,
	//	Page:     0,
	//})
	//
	//if err != nil {
	//	ctx.JSON(http.StatusBadRequest, errorResponse(err))
	//	return
	//}
	//
	//branchUsers := make([]branchUser, 0)
	//
	//for _, sm := range *smList {
	//	branchUsers = append(branchUsers, branchUser{
	//		ID:       int32(sm.UserId),
	//		Avatar:   sm.GetAvatarPointer(),
	//		FullName: sm.FirstName + " " + sm.LastName,
	//	})
	//}
	//
	//ctx.JSON(http.StatusOK, branchUsersResponse{
	//	Result: branchUsers,
	//})
}
