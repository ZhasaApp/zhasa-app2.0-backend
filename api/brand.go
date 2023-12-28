package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/base"
	"zhasa2.0/rating"
	"zhasa2.0/statistic"
)

type SetOwnerDashboardGoalRequest struct {
	BrandId    int32 `json:"brand_id"`
	Value      int64 `json:"value"`
	Month      int32 `json:"month"`
	Year       int32 `json:"year"`
	SaleTypeID int32 `json:"sale_type_id"`
}

func (server *Server) SetOwnerDashboardGoal(ctx *gin.Context) {
	var request SetOwnerDashboardGoalRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	period := statistic.MonthPeriod{
		MonthNumber: request.Month,
		Year:        request.Year,
	}

	err := server.setBrandSaleTypeGoal(request.BrandId, request.SaleTypeID, request.Value, period)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

type GetOwnerDashboardBySaleTypesRequest struct {
	Month   int32 `form:"month" json:"month"`
	Year    int32 `form:"year" json:"year"`
	BrandId int32 `form:"brand_id" json:"brand_id"`
}

type SaleTypeResp struct {
	Title     string `json:"title"`
	Color     string `json:"color"`
	ValueType string `json:"value_type"`
}

type OwnerDashboardBySaleTypesItem struct {
	SaleType SaleTypeResp `json:"sale_type"`
	Achieved int64        `json:"achieved"`
	Goal     int64        `json:"goal"`
}

func (server *Server) GetOwnerDashboardBySaleTypes(ctx *gin.Context) {
	var request GetOwnerDashboardBySaleTypesRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	period := statistic.MonthPeriod{
		MonthNumber: request.Month,
		Year:        request.Year,
	}

	saleTypes, err := server.saleTypeRepo.GetSaleTypes()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result := make([]OwnerDashboardBySaleTypesItem, 0)

	for _, saleType := range *saleTypes {
		achieved, _ := server.getBrandSaleSumFunc(request.BrandId, saleType.Id, period)
		goal, _ := server.getBrandOverallGoalFunc(request.BrandId, saleType.Id, period)

		if goal != 0 {
			result = append(result, OwnerDashboardBySaleTypesItem{
				SaleType: SaleTypeResp{
					Title:     saleType.Title,
					Color:     saleType.Color,
					ValueType: saleType.ValueType,
				},
				Achieved: achieved,
				Goal:     goal,
			})
		}
	}

	ctx.JSON(http.StatusOK, base.ArrayResponse[OwnerDashboardBySaleTypesItem]{
		Result: result,
	})
}

type GetOwnerDashboardByBranchesRequest struct {
	BranchIDs []int32 `form:"branch_ids" json:"branch_ids"`
	Month     int32   `form:"month" json:"month"`
	Year      int32   `form:"year" json:"year"`
	BrandId   int32   `form:"brand_id" json:"brand_id"`
}

type BranchResp struct {
	ID    int32  `json:"id"`
	Title string `json:"title"`
}

type OwnerDashboardByBranchesItem struct {
	Branch      BranchResp                      `json:"branch"`
	SuccessRate float32                         `json:"success_rate"`
	Items       []OwnerDashboardBySaleTypesItem `json:"items"`
}

func (server *Server) GetOwnerDashboardByBranches(ctx *gin.Context) {
	var request GetOwnerDashboardByBranchesRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	period := statistic.MonthPeriod{
		MonthNumber: request.Month,
		Year:        request.Year,
	}

	saleTypes, err := server.saleTypeRepo.GetSaleTypes()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result := make([]OwnerDashboardByBranchesItem, 0)

	ratedBranches, err := server.ratedBranchesFunc(request.BrandId, period, *saleTypes)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	for _, branch := range ratedBranches {
		if len(request.BranchIDs) > 0 && !base.Contains(request.BranchIDs, branch.BranchId) {
			continue
		}
		branchRatioRows := make([]rating.RatioRow, 0)
		items := make([]OwnerDashboardBySaleTypesItem, 0)
		for _, saleType := range *saleTypes {
			salesSum, _ := server.getBranchBrandSaleSumFunc(branch.BranchId, request.BrandId, saleType.Id, period)

			goal, _ := server.getBranchBrandGoalFunc(branch.BranchId, request.BrandId, saleType.Id, period)

			if goal != 0 {
				branchRatioRows = append(branchRatioRows, rating.RatioRow{
					Achieved: salesSum,
					Goal:     goal,
					Gravity:  saleType.Gravity,
				})
				items = append(items, OwnerDashboardBySaleTypesItem{
					SaleType: SaleTypeResp{
						Title:     saleType.Title,
						Color:     saleType.Color,
						ValueType: saleType.ValueType,
					},
					Achieved: salesSum,
					Goal:     goal,
				})
			}
		}
		result = append(result, OwnerDashboardByBranchesItem{
			Branch:      BranchResp{ID: branch.BranchId, Title: branch.Title},
			SuccessRate: rating.CalculateRatio(branchRatioRows) * 100,
			Items:       items,
		})
	}

	ctx.JSON(http.StatusOK, base.ArrayResponse[OwnerDashboardByBranchesItem]{
		Result: result,
	})
}

type GetOwnerGoalRequest struct {
	Month      int32 `form:"month" json:"month"`
	Year       int32 `form:"year" json:"year"`
	BrandId    int32 `form:"brand_id" json:"brand_id"`
	SaleTypeId int32 `form:"sale_type_id" json:"sale_type_id"`
}

type GetOwnerGoalResponse struct {
	Value int64 `json:"value"`
}

func (server *Server) GetOwnerGoal(ctx *gin.Context) {
	var request GetOwnerGoalRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	period := statistic.MonthPeriod{
		MonthNumber: request.Month,
		Year:        request.Year,
	}

	goal, err := server.getBrandOverallGoalFunc(request.BrandId, request.SaleTypeId, period)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, GetOwnerGoalResponse{Value: goal})
}
